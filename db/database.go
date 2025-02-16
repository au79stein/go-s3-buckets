package db

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "uploads.db"

// InitializeDB creates the uploads table if it doesn't exist.
func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	query := `CREATE TABLE IF NOT EXISTS uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_name TEXT,
		s3_prefix TEXT,
		s3_key TEXT UNIQUE,
		file_hash TEXT UNIQUE,
		timestamp TEXT,
		allowed_users TEXT
	)`

	_, err = db.Exec(query)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

// StoreMetadata inserts file upload details into the database.
func StoreMetadata(db *sql.DB, fileName, s3Prefix, s3Key, fileHash string) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := db.Exec(`INSERT INTO uploads (file_name, s3_prefix, s3_key, file_hash, timestamp, allowed_users)
		VALUES (?, ?, ?, ?, ?, ?)`, fileName, s3Prefix, s3Key, fileHash, timestamp, "")
	if err != nil {
		return fmt.Errorf("failed to store metadata: %v", err)
	}
	return nil
}

// CheckFileExistsByHash verifies if a file with the given hash exists in the database.
func CheckFileExistsByHash(db *sql.DB, fileHash string) (string, bool, error) {
	var s3Key string
	err := db.QueryRow("SELECT s3_key FROM uploads WHERE file_hash = ?", fileHash).Scan(&s3Key)
	if err == sql.ErrNoRows {
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("query error: %v", err)
	}
	return s3Key, true, nil
}

// MarkFileAsDeleted updates the database entry for a deleted file.
func MarkFileAsDeleted(db *sql.DB, s3Key string) error {
	_, err := db.Exec("DELETE FROM uploads WHERE s3_key = ?", s3Key)
	if err != nil {
		return fmt.Errorf("failed to mark file as deleted: %v", err)
	}
	return nil
}

