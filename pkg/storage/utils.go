package storage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getFilesFromFolder(folderPath string) ([]string, error) {
	var filePaths []string

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePaths = append(filePaths, filepath.Join(folderPath, entry.Name()))
		}
	}

	return filePaths, nil
}

func getTotalVersionedSize(bucket string, client *s3.Client) (int64, error) {
	var totalSize int64
	paginator := s3.NewListObjectVersionsPaginator(client, &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return 0, err
		}

		for _, version := range page.Versions {
			totalSize += *version.Size
		}
	}

	return totalSize, nil
}

// CreateS3Client returns a MinIO-compatible S3 client
func CreateS3Client(endpoint, accessKey, secretKey string) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "eu-west-4", // default region used by MinIO
			}, nil
		})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-4"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
