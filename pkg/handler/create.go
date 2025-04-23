package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/georgeroumieh-pon/go-client/pkg/client"
	"github.com/minio/minio-go/v7"
)

const (
	teamName string = "cs-team"
)

func CreateBucket(minioClient client.MinioClient) (bucketName string) {
	// Generate bucket name
	today := time.Now().Format("20060102")
	bucketName = fmt.Sprintf("%s-%s", teamName, today)

	// Create the bucket if it doesn't exist
	exists, err := minioClient.Client.BucketExists(minioClient.Ctx, bucketName)
	if err != nil {
		log.Fatalf("Error checking bucket: %v", err)
	}

	if exists {
		fmt.Printf("⚠️ Bucket %s already exists\n", bucketName)
	} else {
		err = minioClient.Client.MakeBucket(minioClient.Ctx, bucketName, minio.MakeBucketOptions{Region: "eu-west-4"})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		fmt.Printf("🪣 Bucket %s created successfully\n", bucketName)
		// Enable versioning
		err = minioClient.Client.EnableVersioning(minioClient.Ctx, bucketName)
		if err != nil {
			log.Fatalf("❌ Failed to enable versioning: %v", err)
		}
		fmt.Printf("🔁 Versioning enabled on bucket %s\n", bucketName)
	}

	return bucketName
}
