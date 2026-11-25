package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
)

// stubStorage is an in-memory fake — no AWS SDK, no network. Used by unit
// tests (inject directly, no env flag needed) and by local dev when
// STUB_MODE=true (wired once in cmd/main.go).
type stubStorage struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewStub constructs the in-memory fake implementation.
func NewStub() Storage {
	return &stubStorage{store: make(map[string][]byte)}
}

func (s *stubStorage) Put(_ context.Context, key string, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = data
	return nil
}

func (s *stubStorage) Get(_ context.Context, key string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.store[key]
	if !ok {
		return nil, fmt.Errorf("stub storage: key %q not found", key)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}
