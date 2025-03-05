package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"my-cucumber-backend/middleware"
	"my-cucumber-backend/models"
	"my-cucumber-backend/services"
)

// GetScenariosHandler retrieves scenarios based on query parameters.
func GetScenariosHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Get optional query parameters
	tagsStr := r.URL.Query().Get("tags")
	folderIDStr := r.URL.Query().Get("folder_id")
	keyword := r.URL.Query().Get("keyword")

	var scenarios []models.Scenario
	userID := user.ID

	// Call the appropriate service function based on the provided parameters
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",") // Split comma-separated tags
		scenarios, err = services.GetScenariosByTags(projectID, userID, tags)
		if err != nil { // Use the outer-scope 'err'
			http.Error(w, "Failed to get scenarios by tags: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if folderIDStr != "" {
		folderID, err := strconv.Atoi(folderIDStr) // Use := to declare a *new* err in this scope
		if err != nil {
			http.Error(w, "Invalid folder ID", http.StatusBadRequest)
			return
		}
		scenarios, err = services.GetScenariosByFolderID(projectID, userID, folderID)
		if err != nil { // Use the newly declared 'err' from this block
			http.Error(w, "Failed to get scenarios by folder ID: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if keyword != "" {
		scenarios, err = services.GetScenariosByName(projectID, userID, keyword)
		if err != nil { // Use the outer-scope 'err'
			http.Error(w, "Failed to get scenarios by name: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		scenarios, err = services.GetScenariosByProjectID(projectID, userID) // Default: get all by project ID
		if err != nil {
			http.Error(w, "Failed to get scenarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Scenarios retrieved successfully",
		"scenarios": scenarios,
	})
}

// RefreshScenariosHandler fetches the latest scenarios from Cucumber Studio and updates the database.
func RefreshScenariosHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	scenarios, err := services.RefreshScenarios(user, projectID) // Call service in services/scenario.go
	if err != nil {
		http.Error(w, "Failed to refresh scenarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Scenarios refreshed successfully",
		"scenarios": scenarios,
	})
}
