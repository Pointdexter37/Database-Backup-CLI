package db

import (
	"context"
)

// Config represents the base connection details for a database
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Database represents the core operations a database dialect must support
type Database interface {
	// Backup initiates the backup process and saves the output to the specified outputPath.
	// It relies on standard CLI tools for each underlying database (like pg_dump).
	Backup(ctx context.Context, config Config, outputPath string) error
}
