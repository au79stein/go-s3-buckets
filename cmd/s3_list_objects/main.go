package main

import (
	"fmt"
	"os"
	"time"

	"gos3listobjects/internal/s3utils"
)

func main() {
	bucket := "cloudnost"
	prefix := ""
	fileType := ".txt" // Example filter
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()
	region := "us-east-1" // Ensure this is dynamic in your implementation

	files, err := s3utils.ListS3Objects(bucket, prefix, fileType, startDate, endDate, region)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

