package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gos3listobjects/db"              // Import SQLite utilities
	"gos3listobjects/internal/s3utils" // Importing custom S3 utility package
)

// computeFileHash computes the SHA-256 hash of a given file.
func computeFileHash(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func main() {
	// Define command-line flags
	bucket := flag.String("bucket", "", "S3 bucket name (required)")
	filePath := flag.String("file", "", "Local file to upload (required)")
	prefix := flag.String("prefix", "", "S3 prefix for file storage")
	region := flag.String("region", "us-east-1", "AWS region")

	flag.Parse()

	// Validate required arguments
	if *bucket == "" || *filePath == "" {
		log.Fatal("Error: --bucket and --file are required")
	}

	// Initialize database connection
	dbConn, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	// Compute file hash
	fileHash, err := computeFileHash(*filePath)
	if err != nil {
		log.Fatalf("Failed to compute file hash: %v", err)
	}

	// Check if the file already exists in the database
	existingS3Key, exists, err := db.CheckFileExistsByHash(dbConn, fileHash)
	if err != nil {
		log.Fatalf("Failed to check existing file: %v", err)
	}
	if exists {
		fmt.Printf("File already exists in S3: %s\n", existingS3Key)
		return
	}

	// Upload file to S3
	s3Key, err := s3utils.UploadFileToS3(*filePath, *bucket, *prefix, *region)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	// DEBUG: Print before storing metadata
	fmt.Printf("DEBUG: Storing in DB -> File: %s | Prefix: %s | S3 Key: %s | Hash: %s\n",
		*filePath, *prefix, s3Key, fileHash)

	// **Store metadata in SQLite**
	err = db.StoreMetadata(dbConn, *filePath, *prefix, s3Key, fileHash)
	if err != nil {
		log.Fatalf("Failed to store metadata in database: %v", err)
	} else {
		fmt.Println("DEBUG: Metadata successfully stored in DB")
	}

	fmt.Printf("Successfully uploaded: %s\n", s3Key)
}

