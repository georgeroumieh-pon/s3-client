package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
	Ctx    context.Context
}
