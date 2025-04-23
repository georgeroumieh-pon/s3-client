package client

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type MinioClient struct {
	Client *minio.Client
	Ctx    context.Context
}
