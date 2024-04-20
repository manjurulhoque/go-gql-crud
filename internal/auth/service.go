package auth

import (
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email     string `json:"email" validate:"required,email,email_exists"`
	Name      string `json:"name" validate:"required"`
	Password1 string `json:"password1" validate:"required,min=5"`
	Password2 string `json:"password2" validate:"required,passwords_match"`
}

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
