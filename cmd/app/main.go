package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/manjurulhoque/go-gql-crud/internal/gql"
	"github.com/manjurulhoque/go-gql-crud/internal/models"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"log/slog"
	"net/http"
)

func main() {
	db, err := dbc.DatabaseConnection()
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())
		panic(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
		slog.Error("Error migrating models")
		return
	}

	r := gin.Default()
	sc := gql.Schema

	// Set up a handler for GraphQL queries
	h := handler.New(&handler.Config{
		Schema:     &sc,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
		//FormatErrorFn: func(err error) gqlerrors.FormattedError {
		//	//if formattedErr, ok := err.(*gqlerrors.Error); ok {
		//	//	return gqlerrors.FormatError(formattedErr)
		//	//}
		//	gqlErr := gqlerrors.FormattedError{
		//		Message: err.Error(),
		//	}
		//	if formattedErr, ok := err.(gqlerrors.FormattedError); ok {
		//		if formattedErr.Extensions != nil {
		//			gqlErr.Extensions = formattedErr.Extensions
		//		}
		//	}
		//	return gqlErr
		//},
	})

	// Convert http.HandlerFunc to gin.HandlerFunc
	gqlHandler := func(c *gin.Context) {
		// Extract the Authorization header from the incoming request
		authHeader := c.GetHeader("Authorization")

		// Create a new context that carries the Authorization header
		ctx := context.WithValue(c.Request.Context(), "Authorization", authHeader)

		// Attach the new context with the Authorization to the request
		c.Request = c.Request.WithContext(ctx)

		// Serve using the GraphQL handler
		h.ServeHTTP(c.Writer, c.Request)
	}

	// Add the GraphQL endpoint
	r.POST("/graphql", gqlHandler)
	r.GET("/graphql", gqlHandler) // Optionally, for tools like GraphiQL

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Error running server", "error", err.Error())
		panic(err)
	}
}
