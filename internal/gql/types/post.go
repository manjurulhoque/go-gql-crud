package types

import (
	"github.com/graphql-go/graphql"
)

var PostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PostType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var PostResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PostResponseType",
	Fields: graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.Boolean,
		},
		"errors": &graphql.Field{
			Type: graphql.NewList(FieldErrorType),
		},
	},
})
