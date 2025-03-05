package models

// Folder represents a folder, including its children for hierarchical representation.
type Folder struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	ParentID *string  `json:"parent_id"` // Use pointer to string to allow for null
	Children []Folder `json:"children"`
}
