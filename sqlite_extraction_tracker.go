package s3utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListS3Objects(bucket, region string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list objects: %w", err)
	}

	var objectKeys []string
	for _, obj := range result.Contents {
		objectKeys = append(objectKeys, *obj.Key)
	}

	return objectKeys, nil
}

