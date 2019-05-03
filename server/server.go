package main

import (
	"database/sql"
	"github.com/santileira/go-graphql/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/santileira/go-graphql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.Connect()
	checkErr(err)
	initDB(db)


	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(go_graphql.NewExecutableSchema(go_graphql.Config{Resolvers: &go_graphql.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB(db *sql.DB) {
	database.MustExec(db,"DROP TABLE IF EXISTS reviews")
	database.MustExec(db,"DROP TABLE IF EXISTS screenshots")
	database.MustExec(db,"DROP TABLE IF EXISTS videos")
	database.MustExec(db,"DROP TABLE IF EXISTS users")
	database.MustExec(db,"CREATE TABLE public.users (id SERIAL PRIMARY KEY, name varchar(255), email varchar(255))")
	database.MustExec(db,"CREATE TABLE public.videos (id SERIAL PRIMARY KEY, name varchar(255), description varchar(255), url text,created_at TIMESTAMP, user_id int, FOREIGN KEY (user_id) REFERENCES users (id))")
	database.MustExec(db,"CREATE TABLE public.screenshots (id SERIAL PRIMARY KEY, video_id int, url text, FOREIGN KEY (video_id) REFERENCES videos (id))")
	database.MustExec(db,"CREATE TABLE public.reviews (id SERIAL PRIMARY KEY, video_id int,user_id int, description varchar(255), rating varchar(255), created_at TIMESTAMP, FOREIGN KEY (user_id) REFERENCES users (id), FOREIGN KEY (video_id) REFERENCES videos (id))")
	database.MustExec(db,"INSERT INTO users(name, email) VALUES('Ridham', 'contact@ridham.me')")
	database.MustExec(db,"INSERT INTO users(name, email) VALUES('Tushar', 'tushar@ridham.me')")
	database.MustExec(db,"INSERT INTO users(name, email) VALUES('Dipen', 'dipen@ridham.me')")
	database.MustExec(db,"INSERT INTO users(name, email) VALUES('Harsh', 'harsh@ridham.me')")
	database.MustExec(db,"INSERT INTO users(name, email) VALUES('Priyank', 'priyank@ridham.me')")
}