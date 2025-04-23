package client

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	endpoint     string = "http://localhost:9000"
	accessKeyEnv string = "ACCESS_KEY"
	secretKeyEnv string = "SECRET_KEY"
)

// CreateS3Client returns a MinIO-compatible S3 client
func CreateS3Client() (s3Client S3Client, err error) {
	ctx := context.Background()
	// Get MinIO credentials from environment variables
	accessKeyID := os.Getenv(accessKeyEnv)
	secretAccessKey := os.Getenv(secretKeyEnv)

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "eu-west-3",
			}, nil
		})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-3"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(awscred.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return s3Client, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	s3Client = S3Client{
		Client: client,
		Ctx:    ctx,
	}

	return s3Client, nil
}
