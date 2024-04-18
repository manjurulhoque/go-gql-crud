package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
)

var PostQueries = graphql.Fields{
	"posts": &graphql.Field{
		Type: graphql.NewList(types.PostType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var posts []models.Post
			err := dbc.GetDB().Find(&posts).Error
			if err != nil {
				return nil, err
			}
			return posts, nil
		},
	},
	"post": &graphql.Field{
		Type: types.PostType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var post models.Post
			err := dbc.GetDB().First(&post, p.Args["id"].(int)).Error
			if err != nil {
				return nil, err
			}
			return post, nil
		},
	},
}
