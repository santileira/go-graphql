package main

import (
	"github.com/rs/cors"
	"github.com/gorilla/websocket"
	"github.com/santileira/go-graphql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query",c.Handler(handler.GraphQL(go_graphql.NewExecutableSchema(go_graphql.Config{Resolvers: &go_graphql.Resolver{}}),

		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}))),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
