package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	MinFileSizeMB      = 10
	MinFileSizeBytes   = MinFileSizeMB * 1024 * 1024
	MaxBucketSizeBytes = 1 * 1024 * 1024 * 1024 // 1 GB
)

// UploadFiles uploads exactly 5 files to the specified bucket in parallel using AWS SDK,
// rejecting files smaller than 10MB and preventing uploads that would exceed 1GB total bucket size.
func UploadFiles(ctx context.Context, s3Client *s3.Client, bucketName string) error {
	filePaths, err := getFilesFromFolder("../files")
	fmt.Println(filePaths)
	if err != nil {
		return fmt.Errorf("❌ Failed to read files folder: %w", err)
	}

	if len(filePaths) < 5 {
		return fmt.Errorf("❌ You must provide minimum 5 files")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(filePaths))
	mu := &sync.Mutex{}

	// Step 1: Get current total bucket size (including all versions)
	currentSize, err := getTotalVersionedSize(bucketName, s3Client)
	if err != nil {
		return fmt.Errorf("❌ Failed to calculate current bucket size: %w", err)
	}

	fmt.Printf("📦 Current bucket size: %d MB\n", currentSize/(1024*1024))

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

			mu.Lock()
			if currentSize+stat.Size() > MaxBucketSizeBytes {
				mu.Unlock()
				errChan <- fmt.Errorf("🚫 Uploading %s would exceed 1GB bucket limit", filePath)
				return
			}
			currentSize += stat.Size()
			mu.Unlock()

			objectKey := filepath.Base(filePath)

			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   file,
			})
			if err != nil {
				errChan <- fmt.Errorf("❌ Upload failed for %s: %w", objectKey, err)
				return
			}
			fmt.Printf("✅ Uploaded %s (%d MB)\n", objectKey, stat.Size()/(1024*1024))

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

// Downloads a file from S3 and saves it to ./download/<objectName>
func DownloadFile(s3Client *s3.Client, bucketName, objectKey string) error {
	ctx := context.TODO()

	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to get object: %w", err)
	}
	defer output.Body.Close()

	// Build the destination path in ./download/
	destPath := filepath.Join("../downloads", filepath.Base(objectKey))

	// Create the destination file
	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to create file: %w", err)
	}
	defer outFile.Close()

	// Copy content
	written, err := io.Copy(outFile, output.Body)
	if err != nil {
		return fmt.Errorf("❌ Failed to write file: %w", err)
	}

	fmt.Printf("✅ Downloaded %s (%d bytes) to %s\n", objectKey, written, destPath)
	return nil
}
