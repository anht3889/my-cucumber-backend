package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"my-cucumber-backend/middleware"
	"my-cucumber-backend/services"
)

// GetFoldersHierarchyHandler retrieves the folder hierarchy for a project.
func GetFoldersHierarchyHandler(w http.ResponseWriter, r *http.Request) {
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
	userID := user.ID

	folders, err := services.GetFoldersHierarchy(projectID, userID)
	if err != nil {
		http.Error(w, "Failed to get folder hierarchy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Folder hierarchy retrieved successfully",
		"folders": folders,
	})
}

// RefreshFoldersHandler fetches the latest folders from Cucumber Studio and updates the database.
func RefreshFoldersHandler(w http.ResponseWriter, r *http.Request) {
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
	//userID, _ := strconv.Atoi(user.ID)

	err = services.RefreshFolders(user, projectID) // Call service in services/folder.go
	if err != nil {
		http.Error(w, "Failed to refresh folders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Folders refreshed successfully",
	})
}
