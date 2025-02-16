package s3utils

import (
	"fmt"
	"time"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ListS3Objects lists objects in an S3 bucket with filtering options.
func ListS3Objects(bucketName, prefix, fileType string, startDate, endDate time.Time, region string) ([]string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	var files []string
	err = svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			// Apply date filtering
			if obj.LastModified.After(startDate) && obj.LastModified.Before(endDate) {
				// Apply file type filtering if needed
				if fileType == "" || (fileType != "" && hasFileType(*obj.Key, fileType)) {
					files = append(files, *obj.Key)
				}
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return files, nil
}

// Helper function to check file type.
func hasFileType(key, fileType string) bool {
	return len(key) >= len(fileType) && key[len(key)-len(fileType):] == fileType
}

// UploadFileToS3 uploads a file to an S3 bucket and returns the uploaded file's key.
func UploadFileToS3(filePath, bucketName, s3Prefix, region string) (string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	s3Key := fmt.Sprintf("%s/%s", s3Prefix, filePath)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return s3Key, nil
}

