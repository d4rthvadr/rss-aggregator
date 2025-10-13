package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/darthvadr/rss-aggregator/internal/database"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	DB *database.Queries
}
func main() {
	log.Println("Starting RSS Aggregator...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("No .env file detected")
	}

	portString := os.Getenv("APP_PORT")
	if portString == "" {
		log.Fatalln("No APP_PORT environment variable detected")
	}

	dbUrlString := os.Getenv("DB_URL")
	if dbUrlString == "" {
		log.Fatalln("No DB_URL environment variable detected")
	}

	db, err := sql.Open("postgres", dbUrlString)
	if err != nil {
		log.Fatalln("Error connecting to database: " + err.Error())
	}
	defer db.Close()


	apiConfig := apiConfig{
		DB: database.New(db),
	}

	log.Println("Listening on port " + portString)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// basic CORS
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the RSS Aggregator!"))
	})
	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.With(apiConfig.middlewareAuth).Get("/users", apiConfig.handlerGetUser)
	v1Router.With(apiConfig.middlewareAuth).Post("/feeds", apiConfig.handlerCreateFeed)
	v1Router.With(apiConfig.middlewareAuth).Get("/feeds", apiConfig.handlerGetFeeds)
	v1Router.With(apiConfig.middlewareAuth).Post("/feed_follows", apiConfig.handlerCreateFeedFollows)


	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Println("Starting server...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting server: " + err.Error())
	}
}
