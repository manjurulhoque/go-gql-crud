package models

import "github.com/manjurulhoque/go-gql-crud/pkg/dbc"

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Create saves a new post to the database
func (post *Post) Create() error {
	return dbc.GetDB().Create(post).Error
}

// Update modifies an existing post in the database
func (post *Post) Update() error {
	return dbc.GetDB().Save(post).Error
}
