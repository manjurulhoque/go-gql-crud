package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/auth"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/internal/utils"
)

type RegisterResponse struct {
	Success bool               `json:"success,omitempty"`
	Errors  []utils.FieldError `json:"errors,omitempty"`
	User    *models.User       `json:"user,omitempty"`
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
}
