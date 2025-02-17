package s3utils

import (
	"fmt"
	//"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Object represents metadata of an S3 object.
type S3Object struct {
	Key          string
	Size         int64
	LastModified time.Time
}

// ListS3Objects lists objects in an S3 bucket with filtering options.
func ListS3Objects(bucketName, prefix, fileType string, startDate, endDate time.Time, region, profile string) ([]S3Object, error) {
	// Set session options
	options := session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	}

	// If a profile is specified, use it
	if profile != "" {
		options.Profile = profile
	}

	// Create AWS session
	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	var objects []S3Object

	// Perform ListObjectsV2 operation
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	for _, obj := range result.Contents {
		// Apply date filtering
		if obj.LastModified.After(startDate) && obj.LastModified.Before(endDate) {
			// Apply file type filtering if needed
			if fileType == "" || hasFileType(*obj.Key, fileType) {
				objects = append(objects, S3Object{
					Key:          *obj.Key,
					Size:         *obj.Size,
					LastModified: *obj.LastModified,
				})
			}
		}
	}

	return objects, nil
}

// Helper function to check file type.
func hasFileType(key, fileType string) bool {
	return len(key) >= len(fileType) && key[len(key)-len(fileType):] == fileType
}

