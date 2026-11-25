package storage

import (
	"context"
	"io"
)

// s3Storage wraps aws-sdk-go-v2/service/s3. TODO: add the real import and
// client construction once --storage s3 is enabled —
//
//	import "github.com/aws/aws-sdk-go-v2/service/s3"
//
// kept dependency-free here so this template compiles without the AWS SDK
// vendored until you actually need it.
type s3Storage struct {
	cfg Config
	// client *s3.Client
}

// NewS3 constructs the real S3/MinIO-backed implementation.
func NewS3(cfg Config) Storage {
	return &s3Storage{cfg: cfg}
}

func (s *s3Storage) Put(ctx context.Context, key string, body io.Reader) error {
	// TODO: s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: &s.cfg.Bucket, Key: &key, Body: body})
	panic("not implemented — wire the real S3 client")
}

func (s *s3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	// TODO: s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.cfg.Bucket, Key: &key})
	panic("not implemented — wire the real S3 client")
}
