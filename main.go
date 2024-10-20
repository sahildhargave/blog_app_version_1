package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv" // Correctly importing the godotenv package
	"github.com/sahildhargave/rss_scraper/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load environment variables from .env file

	godotenv.Load()
	// Print initial message
	// fmt.Println("hello world")

	// Get the PORT environment variable
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn ,err := sql.Open("postgres", dbURL)
	if err!=nil {
		log.Fatal("Can't connect to database", err)
	}

	// queries, err := database.New(conn)
	// if err != nil {
	// 	log.Fatal("Can't connect to db connection:", err)
	// }

	apiCfg := apiConfig{
		DB: database.New(conn),
	}


	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum age of cached preflight responses in seconds. Default is 0 (no cache)
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
    v1Router.Post("/users", apiCfg.handlerUsersCreate)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	// Print the PORT value

}
