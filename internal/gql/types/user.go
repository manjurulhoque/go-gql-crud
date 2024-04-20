package types

import (
	"github.com/graphql-go/graphql"
)

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "UserType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var FieldErrorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FieldError",
	Fields: graphql.Fields{
		"key": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var RegisterResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RegisterResponse",
	Fields: graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.Boolean,
		},
		"errors": &graphql.Field{
			Type: graphql.NewList(FieldErrorType),
		},
		"user": &graphql.Field{
			Type: UserType,
		},
	},
})
