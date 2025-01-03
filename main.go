package main

import (
	"fmt"
	"http-server/api"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Ensure all required variables are present
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" {
		log.Fatal("Database credentials are not fully set in the .env file")
	}

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)

	// Create the server instance
	srv, err := api.NewServer(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start the HTTP server
	port := "8080"
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
