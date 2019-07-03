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
	"os"
)

const defaultPort = "8080"

func main() {

	router := chi.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	config := go_graphql.Config{Resolvers: &go_graphql.Resolver{}}

	config.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []go_graphql.Role) (interface{}, error) {
		roleForContext := auth.ForContext(ctx)

		if roleForContext == nil {
			return nil, errors.UnauthorizedError
		}

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

	router.Use(auth.Middleware())
	router.Use(userdataloader.Middleware())

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(handler.GraphQL(go_graphql.NewExecutableSchema(config),
		handler.ComplexityLimit(10),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}))),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
