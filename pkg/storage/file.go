package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

const (
	minFileSizeMB      = 10
	minFileSizeBytes   = minFileSizeMB * 1024 * 1024 // 10 MB
	maxBucketSizeBytes = 1 * 1024 * 1024 * 1024      // 1 GB
	filesFolderPath    = "../files"
	downloadFolderPath = "../downloads"
)

// UploadFiles uploads files from the local folder "files" to the S3 bucket.
// It apply three rules:
// 1. The files must be at least 5 files and each file must be at least 10MB.
// 2. The total size of the bucket (including all versions) must not exceed 1GB.
// 3. The files must be uploaded concurrently.
func UploadFiles(log *zap.Logger, ctx context.Context, s3Client *s3.Client, bucketName string) (err error) {
	// Get the list of files from the local folder
	filePaths, err := getFilesFromFolder(filesFolderPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read files folder: %w", err)
	}
	// Check if there are at least 5 files
	if len(filePaths) < 5 {
		return fmt.Errorf("❌ You must provide minimum 5 files")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(filePaths))
	mu := &sync.Mutex{}

	// Get current total bucket size (including all versions)
	currentSize, err := getTotalVersionedSize(bucketName, s3Client)
	if err != nil {
		return fmt.Errorf("❌ Failed to calculate current bucket size: %w", err)
	}

	log.Sugar().Infof("Current bucket size: %d MB", currentSize/(1024*1024))

	doneFile6 := make(chan struct{})
	// Iterate over the files and upload them concurrently
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

			if filePath == "../files/file3.json" {
				log.Sugar().Info("Waiting for file6 to complete before uploading file3...")
				<-doneFile6
			}

			// Check if the file is smaller than 10MB
			if stat.Size() < minFileSizeBytes {
				errChan <- fmt.Errorf("⚠️ File %s is smaller than %dMB", filePath, minFileSizeMB)
				return
			}

			mu.Lock()

			objectKey := filepath.Base(filePath)
			bucketSizeAfterUpload := currentSize + stat.Size()

			// Check if the bucket size after upload would exceed the 1GB limit
			if bucketSizeAfterUpload > maxBucketSizeBytes {
				mu.Unlock()
				errChan <- fmt.Errorf("🚫 Uploading %s would exceed 1GB bucket limit", objectKey)
				return
			}
			currentSize += stat.Size()
			mu.Unlock()

			// Upload the file to S3
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   file,
			})
			if err != nil {
				errChan <- fmt.Errorf("❌ Upload failed for %s: %w", objectKey, err)
				return
			}
			log.Sugar().Infof("✅ Uploaded %s (%.2f MB)", objectKey, float64(stat.Size())/(1024*1024))
			if objectKey == "file6.json" {
				log.Sugar().Info("File6 upload completed, notifying file3...")
				close(doneFile6)
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	// Check for errors in the error channel
	var combinedErr error
	for err := range errChan {
		log.Sugar().Info(err)
		combinedErr = fmt.Errorf("some uploads failed")
	}
	return combinedErr
}

// Downloads a file from S3 and saves it to ./download/<objectName>
func DownloadFile(log *zap.Logger, s3Client *s3.Client, bucketName string, filesToDownload []string) error {
	ctx := context.TODO()

	for _, objectKey := range filesToDownload {
		// get the object from S3
		output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			return fmt.Errorf("❌ Failed to get object: %w", err)
		}
		defer output.Body.Close()

		// Build the destination path in ./download/
		destPath := filepath.Join(downloadFolderPath, filepath.Base(objectKey))

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

		log.Sugar().Infof("✅ Downloaded %s (%d bytes) to %s\n", objectKey, written, destPath)
	}
	return nil
}
