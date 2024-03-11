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
			post := models.Post{
				Title:       p.Args["title"].(string),
				Description: p.Args["description"].(string),
			}
			err := post.Create()
			return post, err
		},
	},
	"updatePost": &graphql.Field{
		Type: types.PostType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var post models.Post
			dbc.GetDB().First(&post, p.Args["id"].(int))

			post.Title = p.Args["title"].(string)
			post.Description = p.Args["description"].(string)
			err := post.Update()
			return post, err
		},
	},
}
