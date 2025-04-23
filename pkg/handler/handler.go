package handler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

// getBucketSize returns the total size in bytes of all objects in the given bucket
func getBucketSize(ctx context.Context, client *minio.Client, bucketName string) (int64, error) {
	var totalSize int64

	for object := range client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true}) {
		if object.Err != nil {
			return 0, object.Err
		}
		totalSize += object.Size
	}
	return totalSize, nil
}

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
