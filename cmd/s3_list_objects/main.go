
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gos3listobjects/internal/s3utils"
)

func main() {
	// Define command-line arguments
	bucket := flag.String("bucket", "", "S3 bucket name (required)")
	prefix := flag.String("prefix", "", "S3 prefix to filter objects")
	fileType := flag.String("file-type", ".txt", "File type filter (e.g., .txt, .csv)")
	region := flag.String("region", "us-east-1", "AWS region (optional, defaults to us-east-1)")
	startDateStr := flag.String("start-date", "2024-01-01", "Start date for filtering (YYYY-MM-DD)")
	endDateStr := flag.String("end-date", "", "End date for filtering (YYYY-MM-DD, defaults to current time)")

	flag.Parse()

	// Validate required arguments
	if *bucket == "" {
		fmt.Println("Error: --bucket is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse date arguments
	startDate, err := time.Parse("2006-01-02", *startDateStr)
	if err != nil {
		fmt.Println("Error parsing start-date:", err)
		os.Exit(1)
	}

	endDate := time.Now()
	if *endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			fmt.Println("Error parsing end-date:", err)
			os.Exit(1)
		}
	}

	// List objects
	files, err := s3utils.ListS3Objects(*bucket, *prefix, *fileType, startDate, endDate, *region)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

