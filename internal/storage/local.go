package storage

import (
	"context"
	"fmt"
	"io"
	"os"
)

// LocalStorage implements the Storage interface for saving files to the local disk.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

// Save reads from the reader and writes to a file on the local disk.
func (l *LocalStorage) Save(ctx context.Context, reader io.Reader, filename string) error {
	// If basePath is provided, we could prepend it to filename. 
	// For now, we'll just use filename as the relative/absolute path directly.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to copy data to local file: %w", err)
	}

	return nil
}
