package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	Email  string `json:"email"`
	UserId int    `json:"user_id"`
}

var Secret = "secretAbcdefgh"

func GenerateJWTToken(claims *JWTClaims, expireTime time.Duration) (string, error) {
	claims.ExpiresAt = time.Now().Add(expireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Secret))
}
