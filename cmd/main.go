package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"itmo-devops-sem1-project-template/internal/datasource"
	"itmo-devops-sem1-project-template/internal/routes"
)

func getEnvOrFail(envName string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		log.Fatalf("Environment variable %s is not set", envName)
	}
	return envValue
}

func main() {
	// Retrieve database and server configuration from environment variables
	dbUser := getEnvOrFail("POSTGRES_USER")
	dbPassword := getEnvOrFail("POSTGRES_PASSWORD")
	dbHost := getEnvOrFail("POSTGRES_HOST")
	dbPort := getEnvOrFail("POSTGRES_PORT")
	dbName := getEnvOrFail("POSTGRES_DB")

	// Construct PostgreSQL connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// Establish connection to the database
	db, err := datasource.DBConnect(connString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Initialize server routes
	serverRouter := routes.CreateRoutes(db)

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", serverRouter))
}
