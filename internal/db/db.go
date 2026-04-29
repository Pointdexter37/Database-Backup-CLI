package db

import (
	"context"
	"io"
)

// Config represents the base connection details for a database
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	DockerContainer string
}

// Database represents the core operations a database dialect must support
type Database interface {
	// Backup initiates the backup process and writes the output to the specified io.Writer.
	// It relies on standard CLI tools for each underlying database (like pg_dump).
	Backup(ctx context.Context, config Config, out io.Writer) error
}
