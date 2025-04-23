package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/georgeroumieh-pon/go-client/pkg/client"

	"github.com/minio/minio-go/v7"
)

const (
	MinFileSizeMB      = 10
	MinFileSizeBytes   = MinFileSizeMB * 1024 * 1024
	MaxBucketSizeBytes = 1 * 1024 * 1024 * 1024 // 1 GB
)

// UploadFilesParallel uploads exactly 5 files to the specified bucket in parallel,
// rejecting files smaller than 10MB and preventing uploads that would exceed 1GB total bucket size.
func UploadFiles(minioClient client.MinioClient, bucketName string) error {

	filePaths, err := getFilesFromFolder("../files")
	if err != nil {
		log.Fatalf("❌ Failed to read files folder: %v", err)
	}

	if len(filePaths) < 5 {
		return fmt.Errorf("❌ You must provide exactly 5 files")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(filePaths))
	mu := &sync.Mutex{} // Protects currentSize

	// Step 1: Get the current bucket size
	currentSize, err := getBucketSize(minioClient.Ctx, minioClient.Client, bucketName)
	if err != nil {
		return fmt.Errorf("❌ Failed to calculate current bucket size: %w", err)
	}
	fmt.Println("Current bucket size:", currentSize/(1024*1024), "MB")
	for _, path := range filePaths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			file, err := os.Open(filePath)
			if err != nil {
				errChan <- fmt.Errorf("🚫 Failed to open %s: %w", filePath, err)
				return
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				errChan <- fmt.Errorf("🚫 Failed to stat %s: %w", filePath, err)
				return
			}

			if stat.Size() < MinFileSizeBytes {
				errChan <- fmt.Errorf("⚠️ File %s is smaller than %dMB", filePath, MinFileSizeMB)
				return
			}

			// Lock and check bucket size
			mu.Lock()
			if currentSize+stat.Size() > MaxBucketSizeBytes {
				mu.Unlock()
				errChan <- fmt.Errorf("🚫 Uploading %s would exceed 1GB bucket limit", filePath)
				return
			}
			currentSize += stat.Size()
			mu.Unlock()

			objectName := filepath.Base(filePath)
			_, err = minioClient.Client.PutObject(minioClient.Ctx, bucketName, objectName, file, stat.Size(), minio.PutObjectOptions{})
			if err != nil {
				errChan <- fmt.Errorf("❌ Upload failed for %s: %w", objectName, err)
				return
			}

			fmt.Printf("✅ Uploaded %s (%d MB)\n", objectName, stat.Size()/(1024*1024))
		}(path)
	}

	wg.Wait()
	close(errChan)

	var combinedErr error
	for err := range errChan {
		fmt.Println(err)
		combinedErr = fmt.Errorf("some uploads failed")
	}
	return combinedErr
}
