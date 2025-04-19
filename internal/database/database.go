package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	log.Println("Initializing the database")
	var err error

	dbPath := "./auth.db"
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Failed to up database: %v\n", err)
		return err
	}

	// Test the connection with timeout
	log.Println("Testing database connection...")
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return err
	}

	// Set connection parameters
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Create users table if it doesn't exist
	log.Println("Creating users table if not exists...")

	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}

	createKeyTable := `
	CREATE TABLE IF NOT EXISTS keys ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createKeyTable)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}

	createFoldersTable := `
	        CREATE TABLE IF NOT EXISTS folders (
            id INTEGER PRIMARY KEY AUTOiNCREMENT,
            user_id INTEGER NOT NULL,
            parent_folder_id INTEGER,
            folder_name TEXT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES user(id),
            FOREIGN KEY (parent_folder_id) REFERENCES folders(id) ON DELETE CASCADE
			);`
	_, err = db.Exec(createFoldersTable)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}

	createFilesTable := `
	        CREATE TABLE IF NOT EXISTS files (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            file_name TEXT NOT NULL,
            folder_id INTEGER NOT NULL,
            bin BLOB NOT NULL,
            size INTEGER NOT NULL,
            created_at TEXT,
            FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id)
			);`
	_, err = db.Exec(createFilesTable)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}

	log.Println("Database initialization complete")
	return nil
}

func generateAPIKey() string {
	// Generate a random API key
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CreateUser(username, passwordHash string) error {
	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction: ", err)
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO users(username, password_hash) VALUES(?, ?)", username, passwordHash)
	if err != nil {
		log.Println("Failed to insert user: ", err)
		return err
	}

	// Get the user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %v", err)
	}

	// Generate and insert API key
	apiKey := generateAPIKey()
	_, err = tx.Exec("INSERT INTO keys(user_id, key) VALUES(?, ?)",
		userID, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create API key: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return err
}

type UserData struct {
	PasswordHash string
	APIKey       string
	UserId       int
	FolderId     int64
}

func GetUser(username string) (UserData, error) {
	var userData UserData
	err := db.QueryRow(`
		SELECT u.id, u.password_hash, k.key 
		FROM users u 
		LEFT JOIN keys k 
		ON u.id = k.user_id
		WHERE u.username = ?`, username).Scan(&userData.PasswordHash, &userData.APIKey, &userData.UserId)
	if err != nil {
		return UserData{}, err
	}

	check_root(userData.UserId, &userData)

	return userData, nil
}

func check_root(id int, userData *UserData) error {
	var folderId int64
	err := db.QueryRow("SELECT id FROM folders WHERE user_id = ? AND parent_folder_id IS NULL", id).Scan(&folderId)

	if err == sql.ErrNoRows {
		// Create root folder if it doesn't exist
		result, err := db.Exec(`
		INSERT INTO 
		folders(user_id,parent_folder_id, folder_name) 
		VALUES(?, ?)`, id, nil, "root")

		if err != nil {
			log.Printf("Error creating root folder: %v", err)
			return fmt.Errorf("error creating root folder: %v", err)
		}
		newId, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting new folder ID: %v", err)
			return fmt.Errorf("error getting new folder ID: %v", err)
		}
		folderId = newId

	} else if err != nil {
		log.Println("Error checking root folder: ", err)
		return fmt.Errorf("error checking root folder: %v", err)
	}
	// Update value of folderId in userData
	userData.FolderId = folderId
	return nil

}

func UsernameExists(username string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		log.Println("Error checking username: ", err)
		return false
	}
	return exists
}
