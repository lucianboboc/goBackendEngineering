package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TestAuthenticator struct {
}

const secret = "test"

func (a *TestAuthenticator) GenerateToken(sub int64, iss, aud string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(exp).Unix(),
		"iat": time.Now().Unix(),
		"nbt": time.Now().Unix(),
		"iss": iss,
		"aud": iss,
	})
	return token.SignedString([]byte(secret))
}
func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
