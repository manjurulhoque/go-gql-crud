package models

import "github.com/manjurulhoque/go-gql-crud/pkg/dbc"

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title,min=5"`
	Description string `json:"description"`
	UserId      int    `json:"user_id"`
	User        User   `json:"user"`
}

// Create saves a new post to the database
func (post *Post) Create() error {
	return dbc.GetDB().Create(post).Error
}

// Update modifies an existing post in the database
func (post *Post) Update() error {
	return dbc.GetDB().Save(post).Error
}

// Delete removes an existing post from the database
func (post *Post) Delete() error {
	return dbc.GetDB().Delete(post).Error
}
