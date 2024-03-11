package main

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	},
)

func main() {
	r := gin.Default()

	// Set up a handler for GraphQL queries
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	// Convert http.HandlerFunc to gin.HandlerFunc
	gqlHandler := func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}

	// Add the GraphQL endpoint
	r.POST("/graphql", gqlHandler)
	r.GET("/graphql", gqlHandler) // Optionally, for tools like GraphiQL

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
