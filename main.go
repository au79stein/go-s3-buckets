package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ListObjects lists objects in an S3 bucket with optional filters.
func ListObjects(bucket, prefix, region, fileType string, startDate, endDate time.Time) {
	// Create AWS session with the specified region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	// Pagination setup
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	err = svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			// Apply date filter if provided
			if !startDate.IsZero() && obj.LastModified.Before(startDate) {
				continue
			}
			if !endDate.IsZero() && obj.LastModified.After(endDate) {
				continue
			}

			// Apply file type filter if provided
			if fileType != "" && !hasFileExtension(*obj.Key, fileType) {
				continue
			}

			fmt.Printf("Key: %s | Size: %d | LastModified: %v\n", *obj.Key, *obj.Size, obj.LastModified)
		}
		return true
	})

	if err != nil {
		log.Fatalf("Failed to list objects: %v", err)
	}
}

// hasFileExtension checks if the S3 object key has the specified file extension.
func hasFileExtension(key, extension string) bool {
	return len(key) > len(extension) && key[len(key)-len(extension):] == extension
}

func main() {
	// Define command-line flags
	bucket := flag.String("bucket", "", "S3 bucket name (required)")
	prefix := flag.String("prefix", "", "S3 object prefix filter (optional)")
	region := flag.String("region", os.Getenv("AWS_REGION"), "AWS region (default: AWS_REGION env variable)")
	fileType := flag.String("file-type", "", "Filter by file extension (e.g., .txt, .jpg) (optional)")
	startDateStr := flag.String("start-date", "", "Filter by start date (YYYY-MM-DD) (optional)")
	endDateStr := flag.String("end-date", "", "Filter by end date (YYYY-MM-DD) (optional)")

	flag.Parse()

	// Validate required flag
	if *bucket == "" {
		log.Fatal("Error: --bucket is required")
	}

	// Parse date filters if provided
	var startDate, endDate time.Time
	var err error
	if *startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", *startDateStr)
		if err != nil {
			log.Fatalf("Invalid start-date format: %v", err)
		}
	}
	if *endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			log.Fatalf("Invalid end-date format: %v", err)
		}
	}

	// Default to "us-west-2" if no region is specified
	if *region == "" {
		*region = "us-west-2"
	}

	fmt.Printf("Listing objects in bucket: %s, region: %s\n", *bucket, *region)
	ListObjects(*bucket, *prefix, *region, *fileType, startDate, endDate)
}

