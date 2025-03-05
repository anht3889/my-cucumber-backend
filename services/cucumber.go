package services

import (
	"encoding/json"
	"fmt"
	"io"
	"my-cucumber-backend/models"
	"net/http"
)

const (
	cucumberStudioBaseURL = "https://studio.cucumber.io/api"
)

// ProjectResponse represents the structure of a single project in the Cucumber Studio API response.
type ProjectResponse struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Name string `json:"name"`
		// Other attributes you might need
	} `json:"attributes"`
}

// ProjectsResponse represents the top-level structure of the /projects response.
type ProjectsResponse struct {
	Data []ProjectResponse `json:"data"`
}

// ScenarioResponse represents a single scenario in the Cucumber Studio API response.
type ScenarioResponse struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Name     string `json:"name"`
		FolderID int    `json:"folder-id"` // Directly get folder-id
		// ... other attributes you might need ...
	} `json:"attributes"`
	Relationships struct {
		Tags struct {
			Data []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"tags"`
	} `json:"relationships"`
}

// ScenariosResponse is the top-level response for a list of scenarios.
type ScenariosResponse struct {
	Data     []ScenarioResponse `json:"data"`
	Included []struct {         // For included resources (like tags)
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"included"`
}

// TagResponse represents a tag in the "included" array.
// We'll define a separate struct for clarity, although it overlaps with the anonymous struct above.
type TagResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

// GetProjects fetches projects from Cucumber Studio for a given user and simplifies the data.
func GetProjects(user *models.User) ([]models.Project, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", cucumberStudioBaseURL+"/projects", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.api+json; version=1")
	req.Header.Add("access-token", user.CucumberAccessToken)
	req.Header.Add("client", user.CucumberClientID)
	req.Header.Add("uid", user.Email)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cucumber Studio API returned an error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var projectsResponse ProjectsResponse
	if err := json.Unmarshal(body, &projectsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	// Simplify the data to just ID and Name
	simplifiedProjects := make([]models.Project, len(projectsResponse.Data))
	for i, p := range projectsResponse.Data {
		simplifiedProjects[i] = models.Project{
			ID:   p.ID,
			Name: p.Attributes.Name,
		}
	}

	return simplifiedProjects, nil
}

// FolderResponse represents a single folder in the API response.
type FolderResponse struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Name     string      `json:"name"`
		ParentID json.Number `json:"parent-id"` // MUST BE *string
		// ... other attributes ...
	} `json:"attributes"`
}

// FoldersResponse is the top-level response structure for a list of folders.
type FoldersResponse struct {
	Data []FolderResponse `json:"data"`
}

// GetFolders fetches folders from Cucumber Studio for a given project.
func GetFolders(user *models.User, projectID int) ([]FolderResponse, error) {
	url := fmt.Sprintf("%s/projects/%d/folders", cucumberStudioBaseURL, projectID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.api+json; version=1")
	req.Header.Add("access-token", user.CucumberAccessToken)
	req.Header.Add("client", user.CucumberClientID)
	req.Header.Add("uid", user.Email)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch folders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read body for error message
		return nil, fmt.Errorf("Cucumber Studio API returned an error: %s, Body: %s", resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var foldersResponse FoldersResponse
	if err := json.Unmarshal(body, &foldersResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v\nResponse: %s", err, string(body))
	}

	return foldersResponse.Data, nil
}

// GetScenarios fetches scenarios and their associated tags from Cucumber Studio.
func GetScenarios(user *models.User, projectID int) ([]models.Scenario, error) {
	url := fmt.Sprintf("%s/projects/%d/scenarios?include=tags", cucumberStudioBaseURL, projectID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.api+json; version=1")
	req.Header.Add("access-token", user.CucumberAccessToken)
	req.Header.Add("client", user.CucumberClientID)
	req.Header.Add("uid", user.Email)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scenarios: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cucumber Studio API returned an error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var scenariosResponse ScenariosResponse
	if err := json.Unmarshal(body, &scenariosResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v\nResponse body: %s", err, string(body)) // Include response body for debugging
	}

	// Create a map to look up tags by ID
	tagMap := make(map[string]models.Tag)
	for _, includedItem := range scenariosResponse.Included {
		if includedItem.Type == "tags" {
			tag := models.Tag{
				ID:    includedItem.ID,
				Key:   includedItem.Attributes.Key,
				Value: includedItem.Attributes.Value,
			}
			tagMap[tag.ID] = tag
		}
	}

	// Build the simplified scenario data
	scenarios := make([]models.Scenario, 0, len(scenariosResponse.Data))
	for _, scenarioData := range scenariosResponse.Data {
		scenario := models.Scenario{
			ID:        scenarioData.ID,
			Name:      scenarioData.Attributes.Name,
			FolderID:  scenarioData.Attributes.FolderID,
			ProjectID: projectID,
			Tags:      make([]models.Tag, 0),
		}

		// Associate tags with the scenario using the tagMap
		for _, tagRelationship := range scenarioData.Relationships.Tags.Data {
			if tag, ok := tagMap[tagRelationship.ID]; ok {
				scenario.Tags = append(scenario.Tags, tag)
			}
		}
		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}
