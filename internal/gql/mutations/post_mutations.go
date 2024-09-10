package mutations

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	middleware "github.com/manjurulhoque/go-gql-crud/internal"
	"github.com/manjurulhoque/go-gql-crud/internal/gql/types"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/internal/utils"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
)

type MyExtendedError struct {
	error
	extensions map[string]interface{}
}

func (e *MyExtendedError) Extensions() map[string]interface{} {
	return e.extensions
}

type PostResponse struct {
	Success bool               `json:"success,omitempty"`
	Errors  []utils.FieldError `json:"errors,omitempty"`
}

var PostMutations = graphql.Fields{
	"createPost": &graphql.Field{
		Type: types.PostResponseType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: middleware.AuthMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
			user, ok := p.Context.Value("currentUser").(*models.User)
			if !ok {
				return nil, errors.New("could not retrieve user from context")
			}
			post := models.Post{
				Title:       p.Args["title"].(string),
				Description: p.Args["description"].(string),
				UserId:      user.ID,
			}
			validationErrors := utils.TranslateError(post)
			if validationErrors != nil {
				return &PostResponse{
					Success: false,
					Errors:  validationErrors,
				}, nil
			}

			err := post.Create()
			if err != nil {
				return &PostResponse{
					Success: false,
					Errors: []utils.FieldError{
						{
							Key:   "unknown",
							Value: err.Error(),
						},
					},
				}, nil
			}
			return &PostResponse{
				Success: true,
			}, nil
		}),
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
	"deletePost": &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(int)
			post := models.Post{}
			// Fetch the post to delete by its ID
			err := dbc.GetDB().First(&post, id).Error
			if err != nil {
				//return false, gqlerrors.FormatError(&gqlerrors.Error{
				//	Message: err.Error(),
				//	OriginalError: gqlerrors.ExtendedError(&MyExtendedError{
				//		extensions: map[string]interface{}{
				//			"code": "SOME_ERROR_CODE",
				//		},
				//	}),
				//})
				return false, gqlerrors.FormattedError{
					Message: err.Error(),
					Extensions: map[string]interface{}{
						"code": "NOT_FOUND",
					},
				}
				//return false, err
			}

			// Delete the post
			err = post.Delete()
			if err != nil {
				return false, err // Error during deletion
			}

			return true, nil // Successfully deleted the post
		},
	},
}
