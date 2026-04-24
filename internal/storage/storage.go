package storage

import (
	"context"
	"io"
)

// Storage defines the interface for saving backup streams to a destination.
type Storage interface {
	// Save reads from the provided io.Reader and writes to the destination
	// with the specified filename or key.
	Save(ctx context.Context, reader io.Reader, filename string) error
}
