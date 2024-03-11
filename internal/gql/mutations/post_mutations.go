package mutations

import (
	"fmt"
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
			newPost := models.Post{
				Title:       p.Args["title"].(string),
				Description: p.Args["description"].(string),
			}
			err := dbc.GetDB().Create(&newPost).Error
			if err != nil {
				return nil, err
			}
			return newPost, nil
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
			id := p.Args["id"].(int)
			post := models.Post{}
			err := dbc.GetDB().First(&post, id).Error
			if err != nil {
				return nil, err
			}
			fmt.Println(p.Args)
			title := p.Args["title"].(string)
			description := p.Args["description"].(string)
			post.Title = title
			post.Description = description
			err = dbc.GetDB().Save(post).Error
			if err != nil {
				return nil, err
			}
			return post, err
		},
	},
}
