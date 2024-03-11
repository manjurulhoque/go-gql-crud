package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/mutations"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/queries"
)

func GetRootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			// Initialize your root mutation fields
			rootFields := graphql.Fields{}

			// Combine product mutations
			for name, field := range mutations.PostMutations {
				rootFields[name] = field
			}

			// Add other domain-specific mutations in a similar way

			return rootFields
		}),
	})
}

func GetRootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			rootFields := graphql.Fields{}
			for name, field := range queries.PostQueries {
				rootFields[name] = field
			}
			return rootFields
		}),
	})
}

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    GetRootQuery(),
		Mutation: GetRootMutation(),
	},
)
