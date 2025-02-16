package main

import (
	"flag"
	"fmt"
	"log"

	"gos3listobjects/db"      // Import SQLite utilities
        "gos3listobjects/internal/s3utils" // Importing custom S3 utility package
)

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

	// Upload file to S3 and track in SQLite
	s3Key, err := s3utils.UploadFileToS3(*filePath, *bucket, *prefix, *region)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	fmt.Printf("Successfully uploaded: %s\n", s3Key)
}

