// Package storage is the worked example of the adapter pattern described in
// document/ARCHITECTURE.md — every SDK-based external dependency (AWS,
// Kafka, Firebase, ...) gets an interface here, a real implementation
// wrapping the SDK, and a stub implementation for tests. Delete this
// package if your service doesn't need S3/MinIO (checklist: --storage no);
// copy its shape for whichever SDK-based dependency you do need.
package storage

import (
	"context"
	"io"
)

// Storage is the interface business logic depends on — never the AWS SDK
// client directly. That's what makes New's real-vs-stub choice possible.
type Storage interface {
	Put(ctx context.Context, key string, body io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}

// Config is the subset of service config this adapter needs.
type Config struct {
	Endpoint  string // MinIO locally; real S3 endpoint in deployed envs
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	Stub      bool // set true via STUB_MODE env or per-dependency config flag
}

// New returns the real S3-backed implementation, or the in-memory stub if
// cfg.Stub is set. Call this once in cmd/main.go — never branch on "am I in
// a test" inside business logic that consumes the Storage interface.
func New(cfg Config) Storage {
	if cfg.Stub {
		return NewStub()
	}
	return NewS3(cfg)
}
