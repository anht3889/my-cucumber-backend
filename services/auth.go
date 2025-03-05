package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"my-cucumber-backend/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user.
func CreateUser(email, password, clientID, accessToken string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:               email,
		PasswordHash:        string(hashedPassword),
		CucumberClientID:    clientID,
		CucumberAccessToken: accessToken,
		Projects:            "[]", // Initialize with an empty JSON array
	}

	// Insert the user into the database
	result, err := DB.Exec(
		"INSERT INTO users (email, password_hash, cucumber_client_id, cucumber_access_token, projects) VALUES (?, ?, ?, ?, ?)",
		user.Email, user.PasswordHash, user.CucumberClientID, user.CucumberAccessToken, user.Projects,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	// Get the last inserted ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last inserted ID: %v", err)
	}
	user.ID = int(userID)

	return user, nil
}

// AuthenticateUser checks the provided email and password.
func AuthenticateUser(email, password string) (*models.User, error) {
	user := &models.User{}
	row := DB.QueryRow(
		"SELECT id, email, password_hash, cucumber_client_id, cucumber_access_token, projects FROM users WHERE email = ?",
		email,
	)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CucumberClientID, &user.CucumberAccessToken, &user.Projects)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	row := DB.QueryRow(
		"SELECT id, email, password_hash, cucumber_client_id, cucumber_access_token, projects FROM users WHERE id = ?",
		userID,
	)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CucumberClientID, &user.CucumberAccessToken, &user.Projects)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	return user, nil
}

// UpdateUserProjects updates the projects associated with a user.
func UpdateUserProjects(user *models.User, projects []models.Project) error {
	projectsJSON, err := json.Marshal(projects)
	if err != nil {
		return fmt.Errorf("failed to marshal projects to JSON: %v", err)
	}

	_, err = DB.Exec(
		"UPDATE users SET projects = ? WHERE id = ?",
		string(projectsJSON), user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user projects: %v", err)
	}

	user.Projects = string(projectsJSON) // Update the user object
	return nil
}

// UpdateCucumberCredentials updates a user's Cucumber Studio credentials
func UpdateCucumberCredentials(userID int, clientID, accessToken string) error {
	_, err := DB.Exec(
		"UPDATE users SET cucumber_client_id = ?, cucumber_access_token = ? WHERE id = ?",
		clientID, accessToken, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update cucumber credentials: %v", err)
	}
	return nil
}
