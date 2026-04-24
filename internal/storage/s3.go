package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage implements the Storage interface for saving files to AWS S3.
type S3Storage struct {
	bucket   string
	region   string
	uploader *manager.Uploader
}

// NewS3Storage creates a new S3Storage instance.
func NewS3Storage(ctx context.Context, bucket, region string) (*S3Storage, error) {
	var cfg aws.Config
	var err error

	if region != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Storage{
		bucket:   bucket,
		region:   region,
		uploader: uploader,
	}, nil
}

// Save reads from the reader and uploads the stream to S3.
func (s *S3Storage) Save(ctx context.Context, reader io.Reader, filename string) error {
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})

	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}
