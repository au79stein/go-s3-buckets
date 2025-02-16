package s3utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Debug function to print credentials
func printAWSCredentials(sess *session.Session) {
	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("Failed to retrieve credentials:", err)
		return
	}

	// Ensure credentials are not empty before printing
	if len(creds.AccessKeyID) == 0 || len(creds.SecretAccessKey) == 0 {
		fmt.Println("Credentials appear to be empty! Check AWS profile or environment variables.")
		return
	}

	fmt.Println("AWS_ACCESS_KEY_ID:", creds.AccessKeyID)
	fmt.Println("AWS_SECRET_ACCESS_KEY:", creds.SecretAccessKey)
	fmt.Println("AWS_SESSION_TOKEN:", creds.SessionToken)
}

// ListS3Objects lists objects in an S3 bucket with filtering options.
func ListS3Objects(bucketName, prefix, fileType string, startDate, endDate time.Time, region string) ([]string, error) {
	// Debug: Check AWS_PROFILE
	awsProfile := os.Getenv("AWS_PROFILE")
	log.Println("AWS_PROFILE environment variable:", awsProfile)

	// Explicitly load credentials using profile
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody), // Enable debug logging
		},
		Profile:           "default", // Explicitly set the profile name
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Print credentials to verify they match Python output
	printAWSCredentials(sess)

	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	var files []string

	// Debugging: Try single `ListObjectsV2` request first
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	for _, obj := range result.Contents {
		log.Printf("Object Found: %s (LastModified: %v)", *obj.Key, *obj.LastModified)

		// Apply date filtering
		if obj.LastModified.After(startDate) && obj.LastModified.Before(endDate) {
			// Apply file type filtering if needed
			if fileType == "" || hasFileType(*obj.Key, fileType) {
				files = append(files, *obj.Key)
			}
		}
	}

	// Debugging: Print final matched files
	log.Println("Filtered S3 Files:", files)

	return files, nil
}

// Helper function to check file type.
func hasFileType(key, fileType string) bool {
	return len(key) >= len(fileType) && key[len(key)-len(fileType):] == fileType
}

// UploadFileToS3 uploads a file to an S3 bucket and returns the uploaded file's key.
func UploadFileToS3(filePath, bucketName, s3Prefix, region string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		//Region: aws.String(region),
	})
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

