package client

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	endpoint     string = "http://localhost:9000"
	accessKeyEnv string = "ACCESS_KEY"
	secretKeyEnv string = "SECRET_KEY"
	region       string = "eu-west-4"
)

// S3Client wraps the S3 client and context
func CreateS3Client() (s3Client *S3Client, err error) {
	ctx := context.Background()
	// Get MinIO credentials from environment variables
	accessKey := os.Getenv(accessKeyEnv)
	secretKey := os.Getenv(secretKeyEnv)
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing ACCESS_KEY or SECRET_KEY environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithBaseEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	s3Client = &S3Client{
		Client: client,
		Ctx:    ctx,
	}

	return s3Client, nil
}
