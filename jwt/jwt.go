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
		config contract.Configuration
		store  contract.JwtStore
		secret string
	}

	audience struct {
		audience string
	}
)

var (
	ErrorExpired = errors.New("token parse error")
)

func New(secret string, config contract.Configuration, store contract.JwtStore) contract.Jwt {
	return &Jwt{config: config, store: store, secret: secret}
}

func newAudience(aud string) contract.JwtAudience {
	return &audience{audience: aud}
}

func (j *Jwt) Create(aud contract.JwtAudience) (*contract.Token, error) {
	expireAt := time.Now().Unix() + int64(j.config.GetInt(`lifetime`))

	id := strings2.Join(`-`, aud.Audience(), strconv.FormatInt(time.Now().UnixNano(), 10))

	// Create the Claims
	claim := &jwt.StandardClaims{
		ExpiresAt: expireAt,
		Issuer:    j.config.GetString(`issuer`),
		Audience:  aud.Audience(),
		Id:        id,
		IssuedAt:  time.Now().Unix(),
		//Subject:  `subject`
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(j.secret))
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

func (j *Jwt) Parse(token string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token2, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims2, ok := token2.Claims.(*jwt.StandardClaims); ok && token2.Valid {
		return claims2, nil
	} else {
		return nil, ErrorExpired
	}
}

func (j *Jwt) Invalidate(token string) error {
	claims, err := j.Parse(token)
	if err != nil {
		return err
	}

	err = j.store.Forget(newAudience(claims.Audience))
	if err != nil {
		return err
	}

	return nil
}

func (j *Jwt) Valid(token string) (bool, error) {
	claims, err := j.Parse(token)

	if err != nil {
		return false, err
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return false, errors.New("token is expired")
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
