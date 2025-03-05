package main

import (
	"log"
	"os"

	"my-cucumber-backend/api"
	"my-cucumber-backend/middleware"
	"my-cucumber-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file (optional)")
	}

	// Get database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db"
	}

	// Initialize database
	if err := services.InitializeDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer services.CloseDB()

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}
	middleware.SetSecretKey(secretKey)

	// Initialize Gin router
	r := gin.Default() // Includes Logger and Recovery middleware

	// CORS configuration
	r.Use(middleware.CORSMiddleware())

	// Public routes
	r.POST("/api/register", api.RegisterHandler)
	r.POST("/api/login", api.LoginHandler)
	r.POST("/api/logout", api.LogoutHandler)

	// Protected routes
	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/data", api.ProtectedHandler)
		protected.POST("/refresh-projects", api.RefreshProjectsHandler)
		protected.GET("/scenarios", api.GetScenariosHandler)
		protected.POST("/refresh-scenarios", api.RefreshScenariosHandler)
		protected.GET("/folders", api.GetFoldersHierarchyHandler)
		protected.POST("/refresh-folders", api.RefreshFoldersHandler)
		protected.POST("/charts", api.CreateChartHandler)
		protected.GET("/charts", api.GetChartsHandler)
		protected.POST("/data-tables", api.CreateDataTableHandler)
		protected.GET("/data-tables", api.GetDataTablesHandler)
		protected.PUT("/update-cucumber-credentials", api.UpdateCucumberCredentialsHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	r.Run(":" + port)
}
