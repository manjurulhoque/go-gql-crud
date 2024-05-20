package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/auth"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/internal/utils"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"golang.org/x/crypto/bcrypt"
)

type RegisterResponse struct {
	Success bool               `json:"success,omitempty"`
	Errors  []utils.FieldError `json:"errors,omitempty"`
	User    *models.User       `json:"user,omitempty"`
}

type LoginResponse struct {
	Success bool               `json:"success,omitempty"`
	Errors  []utils.FieldError `json:"errors,omitempty"`
	Access  string             `json:"access"`
	Refresh string             `json:"refresh"`
}

type RefreshAccessResponse struct {
	Success bool               `json:"success,omitempty"`
	Errors  []utils.FieldError `json:"errors,omitempty"`
	Access  string             `json:"access"`
}

var UserMutations = graphql.Fields{
	"register": &graphql.Field{
		Type: types.RegisterResponseType,
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password1": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password2": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			registerInput := auth.RegisterInput{
				Email:     p.Args["email"].(string),
				Name:      p.Args["name"].(string),
				Password1: p.Args["password1"].(string),
				Password2: p.Args["password2"].(string),
			}
			validationErrors := utils.TranslateError(registerInput)
			if validationErrors != nil {
				return &RegisterResponse{
					Success: false,
					Errors:  validationErrors,
					User:    nil,
				}, nil
			}

			err := auth.Register(&registerInput)
			if err != nil {
				return &RegisterResponse{
					Success: false,
				}, err
			}
			return &RegisterResponse{
				Success: true,
				User:    nil,
			}, nil
		},
	},
	"login": &graphql.Field{
		Type: types.LoginResponseType,
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			loginInput := auth.LoginInput{
				Email:    p.Args["email"].(string),
				Password: p.Args["password"].(string),
			}
			validationErrors := utils.TranslateError(loginInput)
			if validationErrors != nil {
				return &LoginResponse{
					Success: false,
					Errors:  validationErrors,
				}, nil
			}
			dbUser := models.User{}
			if err := dbc.GetDB().Table("users").Where("email = ?", loginInput.Email).First(&dbUser).Error; err != nil {
				return &LoginResponse{
					Success: false,
					Errors: []utils.FieldError{
						{
							Key:   "email",
							Value: "Invalid email or password",
						},
					},
				}, err
			}
			err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginInput.Password))

			if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				return &LoginResponse{
					Success: false,
					Errors: []utils.FieldError{
						{
							Key:   "password",
							Value: "Invalid password",
						},
					},
				}, nil
			}

			accessToken, refreshToken, err := auth.Login(&dbUser)

			if err != nil {
				return &LoginResponse{
					Success: false,
					Errors: []utils.FieldError{
						{
							Key:   "unknown",
							Value: err.Error(),
						},
					},
				}, err
			}

			return &LoginResponse{
				Success: true,
				Access:  accessToken,
				Refresh: refreshToken,
			}, nil
		},
	},
	"refreshAccessToken": &graphql.Field{
		Type:        types.TokenRefreshResponseType,
		Description: "Refresh the access token using a refresh token",
		Args: graphql.FieldConfigArgument{
			"refreshToken": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			refreshToken := p.Args["refreshToken"].(string)

			// Validate the refresh token and issue a new access token
			newAccessToken, err := auth.RefreshAccessToken(refreshToken)
			if err != nil {
				return &RefreshAccessResponse{
					Success: false,
					Errors: []utils.FieldError{
						{
							Key:   "unknown",
							Value: err.Error(),
						},
					},
				}, nil
			}

			return &RefreshAccessResponse{
				Success: true,
				Access:  newAccessToken,
			}, nil
		},
	},
}
