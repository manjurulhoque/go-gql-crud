package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
)

var PostMutations = graphql.Fields{
	"createPost": &graphql.Field{
		Type: types.PostType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			newTodo := models.Post{
				Title:       p.Args["title"].(string),
				Description: p.Args["description"].(string),
			}
			err := dbc.GetDB().Create(&newTodo).Error
			if err != nil {
				return nil, err
			}
			return newTodo, nil
		},
	},
}
