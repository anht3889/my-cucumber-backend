package api

import (
	"strconv"
	"strings"

	"my-cucumber-backend/models"
	"my-cucumber-backend/services"

	"github.com/gin-gonic/gin"
)

// GetScenariosHandler retrieves scenarios based on query parameters.
func GetScenariosHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(400, gin.H{"error": "project_id is required"})
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	typedUser := user.(*models.User)
	tagsStr := c.Query("tags")
	folderIDStr := c.Query("folder_id")
	keyword := c.Query("keyword")

	var scenarios []models.Scenario
	userID := typedUser.ID

	// Call the appropriate service function based on the provided parameters
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",") // Split comma-separated tags
		// Each tag should be in format "key:value"
		scenarios, err = services.GetScenariosByTags(projectID, userID, tags)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get scenarios by tags: " + err.Error()})
			return
		}
	} else if folderIDStr != "" {
		folderID, err := strconv.Atoi(folderIDStr) // Use := to declare a *new* err in this scope
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid folder ID"})
			return
		}
		scenarios, err = services.GetScenariosByFolderID(projectID, userID, folderID)
		if err != nil { // Use the newly declared 'err' from this block
			c.JSON(500, gin.H{"error": "Failed to get scenarios by folder ID: " + err.Error()})
			return
		}
	} else if keyword != "" {
		scenarios, err = services.GetScenariosByName(projectID, userID, keyword)
		if err != nil { // Use the outer-scope 'err'
			c.JSON(500, gin.H{"error": "Failed to get scenarios by name: " + err.Error()})
			return
		}
	} else {
		scenarios, err = services.GetScenariosByProjectID(projectID, userID) // Default: get all by project ID
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get scenarios: " + err.Error()})
			return
		}
	}

	c.JSON(200, scenarios)
}

// RefreshScenariosHandler fetches the latest scenarios from Cucumber Studio and updates the database.
func RefreshScenariosHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	typedUser := user.(*models.User)
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(400, gin.H{"error": "project_id is required"})
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	scenarios, err := services.RefreshScenarios(typedUser, projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, scenarios)
}
