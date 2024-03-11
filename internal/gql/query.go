package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
)

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"posts": &graphql.Field{
			Type: graphql.NewList(types.PostType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return []models.Post{
					{
						ID:          1,
						Title:       "Post title",
						Description: "Post description",
					},
				}, nil
			},
		},
	},
})
