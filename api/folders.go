package api

import (
	"my-cucumber-backend/models"
	"my-cucumber-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFoldersHierarchyHandler retrieves the folder hierarchy for a project.
func GetFoldersHierarchyHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectID, err := strconv.Atoi(c.Query("project_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	typedUser := user.(*models.User)
	folders, err := services.GetFoldersHierarchy(projectID, typedUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, folders)
}

// RefreshFoldersHandler fetches the latest folders from Cucumber Studio and updates the database.
func RefreshFoldersHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectID, err := strconv.Atoi(c.Query("project_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	typedUser := user.(*models.User)
	err = services.RefreshFolders(typedUser, projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Folders refreshed successfully"})
}
