package services

import (
	"database/sql"
	"fmt"
	"log"

	"my-cucumber-backend/models"
)

// CreateFolder inserts a new folder into the database.
func CreateFolder(folder *models.Folder, projectID, userID int) error {
	var parentID *string // Use a pointer for nullable parent_id
	if folder.ParentID != nil {
		parentID = folder.ParentID
	}
	_, err := DB.Exec(
		"INSERT INTO folders (id, name, parent_id, project_id, user_id) VALUES (?, ?, ?, ?, ?)",
		folder.ID, folder.Name, parentID, projectID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert folder: %v", err)
	}
	return nil
}

// RefreshFolders fetches folders from Cucumber Studio, deletes existing folders for the project,
// and inserts the new folders.
func RefreshFolders(user *models.User, projectID int) error {
	log.Printf("Refreshing folders for project ID: %d, user ID: %s", projectID, user.ID)

	// 1. Fetch latest folders from Cucumber Studio
	folders, err := GetFolders(user, projectID)
	if err != nil {
		log.Printf("Error fetching folders from Cucumber Studio: %v", err)         // Log the specific error
		return fmt.Errorf("failed to fetch folders from Cucumber Studio: %w", err) // Wrap the error
	}
	log.Println("Successfully fetched folders from Cucumber Studio")

	userID := user.ID

	// 2. Delete existing folders for this project and user
	err = DeleteFoldersByProjectID(projectID, userID)
	if err != nil {
		log.Printf("Error deleting existing folders: %v", err) // Log the specific error
		// Consider not returning here; log the error but try to insert new ones.
		// return fmt.Errorf("failed to delete existing folders: %w", err)
	}
	log.Println("Successfully deleted existing folders (if any)")

	// 3. Insert the new folders.
	for _, folderData := range folders {
		var parentIDPtr *string
		if folderData.Attributes.ParentID != "" {
			parentIDStr := string(folderData.Attributes.ParentID)
			parentIDPtr = &parentIDStr
		}

		folder := models.Folder{
			ID:       folderData.ID,
			Name:     folderData.Attributes.Name,
			ParentID: parentIDPtr,
		}

		err = CreateFolder(&folder, projectID, userID)
		if err != nil {
			log.Printf("Error creating folder (ID: %s): %v", folder.ID, err)          // Log the specific error
			return fmt.Errorf("failed to create folder (ID: %s): %w", folder.ID, err) // Wrap error
		}
		log.Printf("Folder (id: %s) created successfully", folder.ID)
	}
	log.Println("Successfully refreshed folders")
	return nil
}

// getFoldersByProjectID retrieves all folders for a given project and user from the *local database*.
func getFoldersByProjectID(projectID, userID int) ([]models.Folder, error) {
	rows, err := DB.Query(
		"SELECT id, name, parent_id FROM folders WHERE project_id = ? AND user_id = ?",
		projectID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query folders: %v", err)
	}
	defer rows.Close()

	folders := make([]models.Folder, 0)
	for rows.Next() {
		var folder models.Folder
		var parentID sql.NullString // Use sql.NullString for nullable parent_id
		if err := rows.Scan(&folder.ID, &folder.Name, &parentID); err != nil {
			return nil, fmt.Errorf("failed to scan folder row: %v", err)
		}
		if parentID.Valid {
			folder.ParentID = &parentID.String // Set the pointer if valid
		}
		folders = append(folders, folder)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	return folders, nil
}

// GetFoldersHierarchy builds the hierarchical folder structure.
func GetFoldersHierarchy(projectID, userID int) ([]models.Folder, error) {
	allFolders, err := getFoldersByProjectID(projectID, userID)
	if err != nil {
		return nil, err
	}

	// Create a map to look up folders by ID.
	folderMap := make(map[string]*models.Folder)
	for i := range allFolders {
		folderMap[allFolders[i].ID] = &allFolders[i] // Use pointers for efficient modification
	}

	// Build the hierarchy.
	var rootFolders []models.Folder
	for i := range allFolders {
		folder := &allFolders[i] // Work with pointers
		if folder.ParentID == nil {
			// Root folder (no parent)
			rootFolders = append(rootFolders, *folder) // Append a *copy*
		} else {
			// Find the parent and add this folder as a child.
			if parent, ok := folderMap[*folder.ParentID]; ok {
				parent.Children = append(parent.Children, *folder) //Append a *copy*
			} else {
				// Handle cases where the parent ID is invalid (optional).
				log.Printf("Warning: Folder %s has invalid parent ID %s", folder.ID, *folder.ParentID)
				// You could choose to treat it as a root folder, or skip it.
				rootFolders = append(rootFolders, *folder) //Treat as a root folder
			}
		}
	}

	return rootFolders, nil
}

// DeleteFoldersByProjectID deletes all folders associated with a project and user.
func DeleteFoldersByProjectID(projectID, userID int) error {
	_, err := DB.Exec("DELETE FROM folders WHERE project_id = ? and user_id = ?", projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete folders: %v", err)
	}
	return nil
}
