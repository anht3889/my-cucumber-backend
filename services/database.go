package services

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var DB *sql.DB // Exported database connection

// InitializeDB initializes the database connection and creates all necessary tables.
func InitializeDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create the users table if it doesn't exist
	createUserTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            cucumber_client_id TEXT,
            cucumber_access_token TEXT,
            projects TEXT
        );
    `
	_, err = DB.Exec(createUserTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	createScenariosTableSQL := `
        CREATE TABLE IF NOT EXISTS scenarios (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            folder_id INTEGER NOT NULL,
            project_id INTEGER NOT NULL,
            tags TEXT,
            user_id INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `

	_, err = DB.Exec(createScenariosTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create scenarios table: %v", err)
	}
	createFoldersTableSQL := `
        CREATE TABLE IF NOT EXISTS folders (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            parent_id TEXT,
            project_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `
	_, err = DB.Exec(createFoldersTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create folders table:%v", err)
	}
	return nil
}

// CloseDB closes the database connection.
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
