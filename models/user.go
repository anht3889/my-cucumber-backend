package models

type User struct {
	ID                  int    `json:"id"`
	Email               string `json:"email"`
	PasswordHash        string `json:"-"`
	CucumberClientID    string `json:"cucumber_client_id"`
	CucumberAccessToken string `json:"cucumber_access_token"`
	Projects            string `json:"projects"` // Store as JSON string
}

// Project represents a simplified project with just ID and Name.
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
