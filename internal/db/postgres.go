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
	var baseCmd string
	var args []string

	if config.DockerContainer != "" {
		baseCmd = "docker"
		args = []string{"exec", "-i"}
		if config.Password != "" {
			args = append(args, "-e", fmt.Sprintf("PGPASSWORD=%s", config.Password))
		}
		args = append(args, config.DockerContainer, "pg_dump")
	} else {
		baseCmd = "pg_dump"
	}

	pgArgs := []string{
		"-h", config.Host,
		"-U", config.User,
		"-d", config.Database,
	}

	if config.Port > 0 {
		pgArgs = append(pgArgs, "-p", fmt.Sprintf("%d", config.Port))
	}

	args = append(args, pgArgs...)

	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Stdout = out

	// Pass the password via the PGPASSWORD environment variable if not using docker
	if config.DockerContainer == "" && config.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", config.Password))
	}

	// Capture stderr for error reporting
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %v", err)
	}

	return nil
}
