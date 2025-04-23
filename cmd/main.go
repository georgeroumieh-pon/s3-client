package main

import (
	"log"

	"github.com/georgeroumieh-pon/go-client/pkg/client"
	"github.com/georgeroumieh-pon/go-client/pkg/handler"
)

func main() {
	client := client.NewMinioClient()
	bucketName := handler.CreateBucket(client)
	err := handler.UploadFiles(client, bucketName)
	if err != nil {
		log.Fatal(err)
	}
}
