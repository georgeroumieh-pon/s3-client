package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

const (
	teamName string = "cs-team"
)

func CreateBucket(log *zap.Logger, ctx context.Context, s3Client *s3.Client) (bucketName string, err error) {
	// Generate bucket name
	today := time.Now().Format("20060102")
	bucketName = fmt.Sprintf("%s-%s", teamName, today)

	// Check if the bucket exists
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)})
	if err == nil {
		return bucketName, fmt.Errorf("⚠️ Bucket %s already exists", bucketName)
	}

	// Create bucket
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuWest3,
		},
	})
	if err != nil {
		return bucketName, fmt.Errorf("❌ Failed to create bucket: %v", err)

	}

	log.Sugar().Infof("🪣 Bucket %s created successfully\n", bucketName)

	// Enable versioning
	_, err = s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})
	if err != nil {
		return bucketName, fmt.Errorf("❌ Failed to enable versioning: %v", err)

	}

	log.Sugar().Infof("🔁 Versioning enabled on bucket %s", bucketName)

	return bucketName, err
}
