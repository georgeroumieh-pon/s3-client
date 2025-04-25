package main

import (
	"os"

	"github.com/georgeroumieh-pon/go-client/pkg/client"
	"github.com/georgeroumieh-pon/go-client/pkg/storage"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	log := zap.New(core, zap.AddCaller())
	// Create a new S3 client
	s3Client, err := client.CreateS3Client()
	if err != nil {
		log.Sugar().Error(err)
		return
	}
	// Create a new bucket
	bucketName, err := storage.CreateBucket(log, s3Client.Ctx, s3Client.Client)
	if err != nil {
		log.Sugar().Error(err)
		return
	}
	// Upload files to the bucket
	err = storage.UploadFiles(log, s3Client.Ctx, s3Client.Client, bucketName)
	if err != nil {
		log.Sugar().Error(err)
	}
	filesToDownload := []string{"object1.txt", "object2.txt", "object3.txt"}
	// Download a file from the bucket
	err = storage.DownloadFile(log, s3Client.Client, bucketName, filesToDownload)
	if err != nil {
		log.Sugar().Error(err)
	}
}
