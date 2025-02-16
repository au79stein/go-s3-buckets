package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gos3listobjects/internal/s3utils"
)

func main() {
	// Define command-line flags
	bucket := flag.String("bucket", "", "S3 bucket name (required)")
	prefix := flag.String("prefix", "", "S3 object prefix (optional)")
	fileType := flag.String("filetype", "", "File extension filter (e.g., .txt)")
	region := flag.String("region", "us-east-1", "AWS region (default: us-east-1)")
	profile := flag.String("profile", "", "AWS profile to use (optional)")

	flag.Parse()

	// Ensure required flag is provided
	if *bucket == "" {
		fmt.Println("Error: --bucket is required")
		flag.Usage()
		os.Exit(1)
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	// Fetch objects from S3
	files, err := s3utils.ListS3Objects(*bucket, *prefix, *fileType, startDate, endDate, *region, *profile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Print retrieved file list
	if len(files) == 0 {
		fmt.Println("No matching objects found.")
	} else {
		fmt.Println("S3 Objects:")
		for _, file := range files {
			fmt.Println(file)
		}
	}
}

