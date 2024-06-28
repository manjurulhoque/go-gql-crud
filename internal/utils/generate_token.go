package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log/slog"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	Email  string `json:"email"`
	UserId int    `json:"user_id"`
}

var AuthSecret = "secretAbcdefgh"

func GenerateJWTToken(claims *JWTClaims, expireTime time.Duration) (string, error) {
	claims.ExpiresAt = time.Now().Add(expireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AuthSecret))
}

func VerifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(AuthSecret), nil
	})

	if err != nil {
		slog.Error("Error in verifying action", "error", err.Error())
		return nil, errors.New("unauthorized")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	if err := token.Claims.Valid(); err != nil {
		slog.Error("Error in verifying action", "error", err.Error())
		return nil, errors.New("unauthorized")
	}
	return claims, nil
}
