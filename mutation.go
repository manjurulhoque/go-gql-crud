package main

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/models"
	"github.com/manjurulhoque/go-gql-crud/types"
)

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createTodo": &graphql.Field{
			Type: types.PostType,
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				newTodo := models.Post{
					ID:          1,
					Title:       p.Args["title"].(string),
					Description: p.Args["description"].(string),
				}
				//types.todos = append(types.todos, newTodo)
				return newTodo, nil
			},
		},
	},
})
