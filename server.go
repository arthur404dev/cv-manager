package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arthur404dev/cv-manager/config"
	"github.com/arthur404dev/cv-manager/db"
	"github.com/arthur404dev/cv-manager/generated"
	"github.com/arthur404dev/cv-manager/resolvers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	env := config.GetEnv()
	port := env.Port
	if port == "" {
		port = config.DEFAULT_PORT
	}

	db, err := db.New(env.DBName)
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
				AllowCredentials: true,
			},
		),
	)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{DB: db}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
