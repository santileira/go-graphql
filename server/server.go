package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"
	"github.com/pressly/chi"
	"github.com/rs/cors"
	"github.com/santileira/go-graphql"
	"github.com/santileira/go-graphql/api/auth"
	"github.com/santileira/go-graphql/api/dataloaders/user"
	"github.com/santileira/go-graphql/api/errors"
	"log"
	"net/http"
)

const defaultPort = "8080"

func main() {

	router := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	config := go_graphql.Config{Resolvers: &go_graphql.Resolver{}}

	// Authentication
	// Verifies if the role in the context is in the list of roles that can perform this action.
	// This verification is very easy and silly, but you can do a very complexity authentication.
	config.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []go_graphql.Role) (interface{}, error) {

		// Gets role for the context.
		roleForContext := auth.ForContext(ctx)

		// Returns error if not exists role for this action.
		if roleForContext == nil {
			return nil, errors.UnauthorizedError
		}

		// If roleForContext is in the list of roles that can perform this action, continues with the next handler, else return error.
		for _, role := range roles {
			if role == *roleForContext {
				return next(ctx)
			}
		}

		return nil, errors.UnauthorizedError
	}
/*
	// Complexity
	countComplexity := func(childComplexity int) int {
		return 1
	}

	config.Complexity.Query.Videos = countComplexity
*/

	// auth middleware gets value from header "Role" and put it in the context
	router.Use(auth.Middleware())
	router.Use(userdataloader.Middleware())

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(handler.GraphQL(go_graphql.NewExecutableSchema(config),
		handler.ComplexityLimit(100),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}))),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, router))
}
