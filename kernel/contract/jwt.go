package contract

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type (
	Jwt interface {
		Create(audience JwtAudience) (*Token, error)

		Parse(token string) (*jwt.StandardClaims, error)

		Refresh(token string) (*Token, error)

		Invalidate(token string) error

		Valid(token string) (bool, error)
	}

	JwtAudience interface {
		Audience() string
	}

	JwtStore interface {
		Has(id string) bool

		Put(id string, audience JwtAudience, lifetime time.Time) error

		Forget(audience JwtAudience) error
	}

	Token struct {
		Lifetime int64  `json:"lifetime"`
		Token    string `json:"token"`
		Type     string `json:"type"`
	}
)
