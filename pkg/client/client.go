package client

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint     string = "localhost:9000"
	accessKeyEnv string = "ACCESS_KEY"
	secretKeyEnv string = "SECRET_KEY"
)

func NewMinioClient() (minioClient MinioClient) {
	useSSL := false
	ctx := context.Background()

	// Get MinIO credentials from environment variables
	accessKeyID := os.Getenv(accessKeyEnv)
	secretAccessKey := os.Getenv(secretKeyEnv)

	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	minioClient = MinioClient{
		Client: client,
		Ctx:    ctx,
	}

	return minioClient
}
