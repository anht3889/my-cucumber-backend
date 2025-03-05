package main

import (
	"log"
	"net/http"
	"os"

	"my-cucumber-backend/api"
	"my-cucumber-backend/middleware"
	"my-cucumber-backend/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // For loading environment variables
)

func main() {
	// Load environment variables from .env file (if it exists)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file (optional)")
	}

	// Get database path from environment variable (or use a default)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db" // Default database file
	}

	// Initialize the database using the database service
	if err := services.InitializeDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer services.CloseDB() // Close the database connection when main exits

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}
	middleware.SetSecretKey(secretKey)

	r := mux.NewRouter()

	r.HandleFunc("/api/register", api.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", api.LoginHandler).Methods("POST")

	protected := r.PathPrefix("/api/protected").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/data", api.ProtectedHandler).Methods("GET")
	protected.HandleFunc("/refresh-projects", api.RefreshProjectsHandler).Methods("POST")
	protected.HandleFunc("/scenarios", api.GetScenariosHandler).Methods("GET")
	protected.HandleFunc("/refresh-scenarios", api.RefreshScenariosHandler).Methods("POST")
	protected.HandleFunc("/folders", api.GetFoldersHierarchyHandler).Methods("GET")     // Get folder hierarchy
	protected.HandleFunc("/refresh-folders", api.RefreshFoldersHandler).Methods("POST") // Refresh folders

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
