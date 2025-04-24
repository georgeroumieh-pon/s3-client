package main

import (
	"os"

	"github.com/georgeroumieh-pon/go-client/pkg/client"
	"github.com/georgeroumieh-pon/go-client/pkg/storage"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func main() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	log := zap.New(core, zap.AddCaller())

	s3Client, err := client.CreateS3Client()
	if err != nil {
		log.Sugar().Error(err)
	}
	bucketName, err := storage.CreateBucket(log, s3Client.Ctx, s3Client.Client)
	if err != nil {
		log.Sugar().Error(err)
	}
	err = storage.UploadFiles(log, s3Client.Ctx, s3Client.Client, bucketName)
	if err != nil {
		log.Sugar().Error(err)
	}
	err = storage.DownloadFile(log, s3Client.Client, bucketName, "object1.txt")
	if err != nil {
		log.Sugar().Error(err)
	}
}
