package models

// Tag represents a simplified tag with just ID and Name.
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"` // Combined key and value for simplicity
}

// Scenario represents a simplified scenario with ID, Name, FolderID, and Tags.
type Scenario struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FolderID int    `json:"folder_id"`
	Tags     []Tag  `json:"tags"`
}
