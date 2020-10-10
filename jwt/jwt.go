package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/firmeve/firmeve/kernel/contract"
	strings2 "github.com/firmeve/firmeve/support/strings"
	"strconv"
	"time"
)

type (
	Jwt struct {
		config *Configuration
		store  contract.JwtStore
	}

	audience struct {
		audience string
	}

	Configuration struct {
		Secret   string `json:"secret" yaml:"secret"`
		Issuer   string `json:"issuer" yaml:"issuer"`
		Lifetime int    `json:"lifetime" yaml:"lifetime"`
	}
)

var (
	ErrorExpired = errors.New("token expired")
)

func New(config *Configuration, store contract.JwtStore) contract.Jwt {
	return &Jwt{config: config, store: store}
}

func newAudience(aud string) contract.JwtAudience {
	return &audience{audience: aud}
}

func (j *Jwt) Create(aud contract.JwtAudience) (*contract.Token, error) {
	expireAt := time.Now().Unix() + int64(j.config.Lifetime)

	id := strings2.Join(`-`, aud.Audience(), strconv.FormatInt(time.Now().UnixNano(), 10))

	// Create the Claims
	claim := &jwt.StandardClaims{
		ExpiresAt: expireAt,
		Issuer:    j.config.Issuer,
		Audience:  aud.Audience(),
		Id:        id,
		IssuedAt:  time.Now().Unix(),
		//Subject:  `subject`
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(j.config.Secret))
	if err != nil {
		return nil, err
	}

	if err := j.store.Put(id, aud, time.Unix(expireAt, 0)); err != nil {
		return nil, err
	}

	return &contract.Token{
		Type:     "Bearer",
		Lifetime: expireAt,
		Token:    token,
	}, err
}

// Just parse without any valid verification
func (j *Jwt) Parse(token string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token2, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	})

	return token2.Claims.(*jwt.StandardClaims), err
}

func (j *Jwt) Invalidate(token string) error {
	claims, _ := j.Parse(token)

	err := j.store.Forget(newAudience(claims.Audience))
	if err != nil {
		return err
	}

	return nil
}

func (j *Jwt) Valid(token string) (bool, error) {
	claims, err := j.Parse(token)

	if err != nil {
		// 过期单独处理
		if e, ok := err.(*jwt.ValidationError); ok && e.Errors == jwt.ValidationErrorExpired {
			return false, ErrorExpired
		}
		return false, err
	}

	if j.store.Has(claims.Id) {
		return true, nil
	}

	return false, errors.New(`token not found`)
}

func (j *Jwt) Refresh(token string) (*contract.Token, error) {
	claims, err := j.Parse(token)
	if err != nil {
		return nil, err
	}

	if err := j.Invalidate(token); err != nil {
		return nil, err
	}

	return j.Create(newAudience(claims.Audience))
}

func (a *audience) Audience() string {
	return a.audience
}
