package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gos3listobjects/internal/s3utils"
)

func main() {
	// Define and parse command-line arguments
	bucket := flag.String("bucket", "", "Name of the S3 bucket (required)")
	flag.Parse()

	if *bucket == "" {
		fmt.Fprintln(os.Stderr, "Error: --bucket argument is required")
		os.Exit(1)
	}

	region := flag.String("region", "", "aws region for bucket")
	flag.Parse()

	if *region == "" {
		*region = "us-east-1"
	}

	prefix := ""
	fileType := ".txt" // Example filter
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	files, err := s3utils.ListS3Objects(*bucket, prefix, fileType, startDate, endDate, *region)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

