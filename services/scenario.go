package services

import (
	"encoding/json"
	"fmt"
	"log"
	"my-cucumber-backend/models"
	"strconv"
	"strings"
)

// CreateScenario creates a scenario record.
func CreateScenario(scenario *models.Scenario, projectID, userID int) error {
	tagsJSON, err := json.Marshal(scenario.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags to JSON: %v", err)
	}

	_, err = DB.Exec(
		"INSERT INTO scenarios (id, name, folder_id, project_id, tags, user_id) VALUES (?, ?, ?, ?, ?, ?)",
		scenario.ID, scenario.Name, scenario.FolderID, projectID, string(tagsJSON), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert scenario: %v", err)
	}
	return nil
}

// GetScenariosByProjectID retrieves all scenarios for a given project and user.
func GetScenariosByProjectID(projectID, userID int) ([]models.Scenario, error) {
	rows, err := DB.Query(
		"SELECT id, name, folder_id, project_id, tags FROM scenarios WHERE project_id = ? AND user_id = ?",
		projectID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenarios: %v", err)
	}
	defer rows.Close()

	scenarios := make([]models.Scenario, 0)
	for rows.Next() {
		var scenario models.Scenario
		var tagsJSON string
		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.FolderID, &scenario.ProjectID, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan scenario row: %v", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &scenario.Tags); err != nil {
			// Log the error, but maybe don't fail the entire operation.
			log.Printf("Error unmarshaling tags for scenario %s: %v", scenario.ID, err)
			scenario.Tags = []models.Tag{} // Set to empty to avoid returning garbage data
		}
		scenarios = append(scenarios, scenario)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	return scenarios, nil
}

// GetScenariosByTags retrieves scenarios matching ALL provided tags for a given project and user.
func GetScenariosByTags(projectID, userID int, tags []string) ([]models.Scenario, error) {
	if len(tags) == 0 {
		return []models.Scenario{}, nil
	}

	// Build the query to match both key and value
	query := `
        SELECT DISTINCT s.id, s.name, s.folder_id, s.project_id, s.tags
        FROM scenarios s
        WHERE s.project_id = ? AND s.user_id = ?
    `
	args := []interface{}{projectID, userID}

	// For each tag in format "key:value", add a condition
	for _, tag := range tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			continue
		}
		// Add JSON path condition to match key and value pair
		query += ` AND s.tags LIKE ?`
		// Look for both key and value in the JSON
		args = append(args, fmt.Sprintf(`%%"key":"%s","value":"%s"%%`, parts[0], parts[1]))
	}

	log.Printf("Query: %s, Args: %v", query, args) // Add logging for debugging

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenarios by tags: %v", err)
	}
	defer rows.Close()

	scenarios := make([]models.Scenario, 0)
	for rows.Next() {
		var scenario models.Scenario
		var tagsJSON string

		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.FolderID, &scenario.ProjectID, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan scenario row: %v", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &scenario.Tags); err != nil {
			log.Printf("Error unmarshaling tags for scenario %s: %v", scenario.ID, err)
			scenario.Tags = []models.Tag{}
		}

		// Filter scenarios to ensure ALL provided tags match
		if matchesAllTags(scenario.Tags, tags) {
			scenarios = append(scenarios, scenario)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return scenarios, nil
}

// matchesAllTags checks if a scenario's tags match all the provided tag filters
func matchesAllTags(scenarioTags []models.Tag, tagFilters []string) bool {
	for _, filter := range tagFilters {
		found := false
		// Split the filter into key and value
		parts := strings.Split(filter, ":")
		if len(parts) != 2 {
			continue // Skip invalid filters
		}
		filterKey := parts[0]
		filterValue := parts[1]

		// Check if any tag matches both key and value
		for _, tag := range scenarioTags {
			if tag.Key == filterKey && tag.Value == filterValue {
				found = true
				break
			}
		}
		if !found {
			return false // If any filter doesn't match, return false
		}
	}
	return true
}

// GetScenariosByFolderID retrieves scenarios within a specific folder.
func GetScenariosByFolderID(projectID, userID, folderID int) ([]models.Scenario, error) {
	rows, err := DB.Query(
		"SELECT id, name, folder_id, project_id, tags FROM scenarios WHERE project_id = ? AND user_id = ? AND folder_id = ?",
		projectID, userID, folderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenarios by folder ID: %v", err)
	}
	defer rows.Close()

	scenarios := make([]models.Scenario, 0)
	for rows.Next() {
		var scenario models.Scenario
		var tagsJSON string
		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.FolderID, &scenario.ProjectID, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan scenario row: %v", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &scenario.Tags); err != nil {
			log.Printf("Error unmarshaling tags for scenario %s: %v", scenario.ID, err)
			scenario.Tags = []models.Tag{}
		}
		scenarios = append(scenarios, scenario)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return scenarios, nil
}

// GetScenariosByName retrieves scenarios containing a keyword in their name.
func GetScenariosByName(projectID, userID int, keyword string) ([]models.Scenario, error) {
	rows, err := DB.Query(
		"SELECT id, name, folder_id, project_id, tags FROM scenarios WHERE project_id = ? AND user_id = ? AND name LIKE ?",
		projectID, userID, "%"+keyword+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenarios by name: %v", err)
	}
	defer rows.Close()

	scenarios := make([]models.Scenario, 0)
	for rows.Next() {
		var scenario models.Scenario
		var tagsJSON string
		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.FolderID, &scenario.ProjectID, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan scenario row: %v", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &scenario.Tags); err != nil {
			log.Printf("Error unmarshaling tags for scenario %s: %v", scenario.ID, err)
			scenario.Tags = []models.Tag{}
		}
		scenarios = append(scenarios, scenario)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	return scenarios, nil
}

// DeleteScenariosByProjectID deletes all scenarios associated with a project and user.
func DeleteScenariosByProjectID(projectID, userID int) error {
	_, err := DB.Exec("DELETE FROM scenarios WHERE project_id = ? and user_id = ?", projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete scenarios: %v", err)
	}
	return nil
}

// RefreshScenarios fetches and updates scenarios from Cucumber Studio.
func RefreshScenarios(user *models.User, projectID int) ([]models.Scenario, error) {
	// 1. Fetch latest scenarios from Cucumber Studio
	scenarios, err := GetScenarios(user, projectID) // Use existing function
	if err != nil {
		return nil, err
	}

	// 2. Delete existing scenarios for this project and user
	err = DeleteScenariosByProjectID(projectID, user.ID) // user.ID is already an int
	if err != nil {
		// Log the error.  Consider whether to continue or abort.
		log.Printf("Failed to delete existing scenarios for project %d, user %d: %v", projectID, user.ID, err)
		// return nil, err  // Option: Abort on deletion failure
	}

	// 3. Insert the new scenarios.
	for _, scenario := range scenarios {
		err = CreateScenario(&scenario, projectID, user.ID) // user.ID is already an int
		if err != nil {
			// Log and decide how to handle individual errors (e.g., continue or abort)
			log.Printf("Failed to create scenario: %v", err)
			return nil, err // Option: Abort if a single scenario fails to insert
		}
	}
	return scenarios, nil
}

// RefreshAllScenarios fetches and updates scenarios from all associated projects.
func RefreshAllScenarios(user *models.User) ([]models.Scenario, error) {
	// 1. Get all associated projects
	projects, err := GetProjects(user)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %v", err)
	}

	// 2. Refresh scenarios for each project
	var allScenarios []models.Scenario
	for _, project := range projects {
		projectID, err := strconv.Atoi(project.ID)
		if err != nil {
			log.Printf("Invalid project ID %s: %v", project.ID, err)
			continue
		}

		scenarios, err := RefreshScenarios(user, projectID)
		if err != nil {
			// Log the error but continue with other projects
			log.Printf("Failed to refresh scenarios for project %d: %v", projectID, err)
			continue
		}

		allScenarios = append(allScenarios, scenarios...)
	}

	if len(allScenarios) == 0 && len(projects) > 0 {
		return nil, fmt.Errorf("failed to refresh scenarios for any project")
	}

	return allScenarios, nil
}
