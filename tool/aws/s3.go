package aws

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
}

func NewS3(config IConfig) *S3 {
	return &S3{
		client: s3.NewFromConfig(config.GetInstance()),
	}
}

func (s *S3) Upload(ctx context.Context, bucketName string, fileName string, fileBytes []byte, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	return err
}
