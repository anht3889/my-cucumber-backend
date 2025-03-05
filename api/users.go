package api

import (
	"encoding/json"
	"time"

	"my-cucumber-backend/middleware"
	"my-cucumber-backend/models"
	"my-cucumber-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RegisterHandler handles user registration.
func RegisterHandler(c *gin.Context) {
	var req struct {
		Email               string `json:"email" binding:"required"`
		Password            string `json:"password" binding:"required"`
		CucumberClientID    string `json:"cucumber_client_id" binding:"required"`
		CucumberAccessToken string `json:"cucumber_access_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := services.CreateUser(req.Email, req.Password, req.CucumberClientID, req.CucumberAccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	projects, err := services.GetProjects(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch initial projects: " + err.Error()})
		return
	}

	if err := services.UpdateUserProjects(user, projects); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user projects: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"projects": projects,
	})
}

// LoginHandler handles user login and issues a JWT.
func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(middleware.SecretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

// RefreshProjectsHandler allows a user to refresh their list of projects.
func RefreshProjectsHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	typedUser := user.(*models.User)
	projects, err := services.GetProjects(typedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch projects: " + err.Error()})
		return
	}

	if err := services.UpdateUserProjects(typedUser, projects); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user projects: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Projects refreshed successfully",
		"projects": projects,
	})
}

// ProtectedHandler is an example of a protected route.
func ProtectedHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	typedUser := user.(*models.User)
	var projects []models.Project
	if err := json.Unmarshal([]byte(typedUser.Projects), &projects); err != nil {
		c.JSON(500, gin.H{"error": "Failed to unmarshal projects"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Hello, here is the protected data.",
		"email":    typedUser.Email,
		"projects": projects,
	})
}

// UpdateCucumberCredentialsHandler handles updating a user's Cucumber Studio credentials
func UpdateCucumberCredentialsHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		CucumberClientID    string `json:"cucumber_client_id" binding:"required"`
		CucumberAccessToken string `json:"cucumber_access_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	typedUser := user.(*models.User)
	if err := services.UpdateCucumberCredentials(typedUser.ID, req.CucumberClientID, req.CucumberAccessToken); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cucumber credentials updated successfully",
	})
}
