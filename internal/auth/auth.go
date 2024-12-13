package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Authenticator interface {
	GenerateToken(sub int64, iss, aud string, exp time.Duration) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
