package models

// Tag represents a tag with ID, Key and Value.
type Tag struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Scenario represents a simplified scenario with ID, Name, FolderID, ProjectID, and Tags.
type Scenario struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FolderID  int    `json:"folder_id"`
	ProjectID int    `json:"project_id"`
	Tags      []Tag  `json:"tags"`
}
