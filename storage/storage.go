package storage

// Note: This package is kept for backwards compatibility
// The SDK works without database storage - all operations are in-memory
// If you need persistent storage, implement the interfaces below

import (
	"context"
	"io"
)

// FileStorage interface for file operations
type FileStorage interface {
	Save(ctx context.Context, filename string, data io.Reader) (string, error)
	Get(ctx context.Context, id string) (io.ReadCloser, error)
	Delete(ctx context.Context, id string) error
}

// InMemoryStorage provides a simple in-memory file storage
type InMemoryStorage struct {
	files map[string][]byte
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		files: make(map[string][]byte),
	}
}

// Save stores file data in memory
func (s *InMemoryStorage) Save(_ context.Context, filename string, data io.Reader) (string, error) {
	content, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}
	s.files[filename] = content
	return filename, nil
}

// Get retrieves file data from memory
func (s *InMemoryStorage) Get(_ context.Context, id string) ([]byte, bool) {
	data, ok := s.files[id]
	return data, ok
}

// Delete removes file from memory
func (s *InMemoryStorage) Delete(_ context.Context, id string) error {
	delete(s.files, id)
	return nil
}
