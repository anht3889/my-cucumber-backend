package api

import (
	"encoding/json"
	"net/http"
	"time"

	"my-cucumber-backend/middleware"
	"my-cucumber-backend/models"
	"my-cucumber-backend/services"

	"github.com/golang-jwt/jwt/v5"
)

// RegisterHandler handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email               string `json:"email"`
		Password            string `json:"password"`
		CucumberClientID    string `json:"cucumber_client_id"`
		CucumberAccessToken string `json:"cucumber_access_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := services.CreateUser(requestBody.Email, requestBody.Password, requestBody.CucumberClientID, requestBody.CucumberAccessToken)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Fetch and store initial projects
	projects, err := services.GetProjects(user)
	if err != nil {
		http.Error(w, "Failed to fetch initial projects: "+err.Error(), http.StatusInternalServerError)
		// Consider:  Do you want to *fail* registration if fetching projects fails?
		//            Or should you still create the user and log the error?
		return
	}

	if err := services.UpdateUserProjects(user, projects); err != nil {
		http.Error(w, "Failed to update user projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the Content-Type
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"projects": projects, // Return the projects in the registration response
	})
}

// LoginHandler handles user login and issues a JWT.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := services.AuthenticateUser(requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Use the secretKey variable directly from the middleware package
	tokenString, err := token.SignedString([]byte(middleware.SecretKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // Set Content-Type
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// RefreshProjectsHandler allows a user to refresh their list of projects.
func RefreshProjectsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projects, err := services.GetProjects(user)
	if err != nil {
		http.Error(w, "Failed to fetch projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := services.UpdateUserProjects(user, projects); err != nil {
		http.Error(w, "Failed to update user projects: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// The context update is NOT needed and has been removed:
	// ctx := r.Context()
	// ctx = context.WithValue(ctx, middleware.UserContextKey, user)
	// r = r.WithContext(ctx)

	w.Header().Set("Content-Type", "application/json") // Set Content-Type
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Projects refreshed successfully",
		"projects": projects,
	})
}

// ProtectedHandler is an example of a protected route.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var projects []models.Project
	err := json.Unmarshal([]byte(user.Projects), &projects)
	if err != nil {
		http.Error(w, "Failed to unmarshal projects", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // Set Content-Type

	// You can now use the user to fetch user-specific data, etc.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Hello, here is the protected data.",
		"email":    user.Email, // Use the user object from the context
		"projects": projects,
	})
}
