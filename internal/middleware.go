package middleware

import (
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/db"
	"github.com/manjurulhoque/go-gql-crud/internal/utils"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"strings"
)

func AuthMiddleware(next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {

		token, ok := p.Context.Value("Authorization").(string)
		if !ok || token == "" {
			return nil, errors.New("access denied: no Authorization header provided")
		}
		bearerToken := ""
		if len(strings.Split(token, " ")) == 2 {
			bearerToken = strings.Split(token, " ")[1]
		}
		if bearerToken == "" {
			return nil, errors.New("access denied: user is not authenticated")
		}

		claims, err := utils.VerifyAction(bearerToken)
		if err != nil {
			return nil, err
		}
		userRepository := db.NewUserRepository(dbc.GetDB())
		user, _ := userRepository.FindUserByEmail(claims.Email)

		if user.Email != claims.Email {
			return nil, errors.New("unauthorized")
		}

		// Create a new context that carries the user object
		ctxWithUser := context.WithValue(p.Context, "currentUser", user)

		// Proceed with the actual resolver, passing the new context
		return next(graphql.ResolveParams{
			Context: ctxWithUser,
			Args:    p.Args,
			Info:    p.Info,
			Source:  p.Source,
		})
	}
}
