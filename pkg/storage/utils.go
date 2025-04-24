package storage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// getFilesFromFolder returns a list of file paths from the specified folder.
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

// getTotalVersionedSize calculates the total size of all versions of objects in the specified S3 bucket.
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
