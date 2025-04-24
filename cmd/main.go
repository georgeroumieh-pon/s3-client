package main

import (
	"log"

	"github.com/georgeroumieh-pon/go-client/pkg/client"
	"github.com/georgeroumieh-pon/go-client/pkg/storage"
)

func main() {
	s3Client, err := client.CreateS3Client()
	if err != nil {
		log.Fatal(err)
	}
	bucketName := storage.CreateBucket(s3Client.Ctx, s3Client.Client)
	err = storage.UploadFiles(s3Client.Ctx, s3Client.Client, bucketName)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.DownloadFile(s3Client.Client, bucketName, "object1.txt")
	if err != nil {
		log.Fatal(err)
	}
}
