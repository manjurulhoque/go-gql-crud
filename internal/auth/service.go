package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/internal/utils"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterInput struct {
	Email     string `json:"email" validate:"required,email,email_exists"`
	Name      string `json:"name" validate:"required"`
	Password1 string `json:"password1" validate:"required,min=5"`
	Password2 string `json:"password2" validate:"required,passwords_match"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

const (
	AccessTokenExpireTime  = 24 * time.Hour     // 1 day
	RefreshTokenExpireTime = 7 * 24 * time.Hour // 7 days
)

func Register(registerInput *RegisterInput) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password1), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Email:    registerInput.Email,
		Name:     registerInput.Name,
		Password: string(hashedPassword),
	}

	if err = dbc.GetDB().Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func Login(user *models.User) (string, string, error) {
	accessClaims := &utils.JWTClaims{
		Email:  user.Email,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}
	refreshClaims := &utils.JWTClaims{
		Email:  user.Email,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}

	accessToken, err := utils.GenerateJWTToken(accessClaims, AccessTokenExpireTime)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateJWTToken(refreshClaims, RefreshTokenExpireTime)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.AuthSecret), nil
	})
	if err != nil {
		return "", errors.New("couldn't parse refresh token")
	}
	if claims, ok := token.Claims.(*utils.JWTClaims); ok && token.Valid {
		if claims.ExpiresAt > time.Now().Unix() {
			return utils.GenerateJWTToken(&utils.JWTClaims{
				UserId: claims.UserId,
				Email:  claims.Email, // Optionally include email if needed
			}, AccessTokenExpireTime)
		}
	}
	return "", errors.New("invalid refresh token")
}
