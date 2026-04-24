package db

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// PostgresDB implements the Database interface for PostgreSQL
type PostgresDB struct{}

// NewPostgresDB creates a new PostgresDB instance
func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

// Backup executes pg_dump to backup a PostgreSQL database
func (p *PostgresDB) Backup(ctx context.Context, config Config, out io.Writer) error {
	args := []string{
		"-h", config.Host,
		"-U", config.User,
		"-d", config.Database,
	}

	if config.Port > 0 {
		args = append(args, "-p", fmt.Sprintf("%d", config.Port))
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	cmd.Stdout = out

	// Pass the password via the PGPASSWORD environment variable
	if config.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", config.Password))
	}

	// Capture stderr for error reporting
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %v", err)
	}

	return nil
}
