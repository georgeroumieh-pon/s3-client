package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	teamName string = "cs-team"
)

func CreateBucket(ctx context.Context, s3Client *s3.Client) (bucketName string) {
	// Generate bucket name
	today := time.Now().Format("20060102")
	bucketName = fmt.Sprintf("%s-%s", teamName, today)

	// Check if the bucket exists
	_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err == nil {
		fmt.Printf("⚠️ Bucket %s already exists\n", bucketName)
		return bucketName
	}

	// Create bucket
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuWest3,
		},
	})

	if err != nil {
		log.Fatalf("❌ Failed to create bucket: %v", err)
	}

	fmt.Printf("🪣 Bucket %s created successfully\n", bucketName)

	// Enable versioning
	_, err = s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})

	if err != nil {
		log.Fatalf("❌ Failed to enable versioning: %v", err)
	}

	fmt.Printf("🔁 Versioning enabled on bucket %s\n", bucketName)

	return bucketName
}
