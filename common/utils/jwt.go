package utils

import (
	"github.com/golang-jwt/jwt/v4"
)

// NewJwtToken 生成一个 jwt token
func NewJwtToken(payload map[string]interface{}, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))
	return token.SignedString([]byte(secret))
}
