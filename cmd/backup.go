package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"dbbackup/internal/db"
	"dbbackup/internal/processor"
	"dbbackup/internal/scheduler"
	"dbbackup/internal/storage"

	"github.com/spf13/cobra"
)

var (
	dbType     string
	dbName     string
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	outputPath string
	compress   bool
	storageDst string
	s3Bucket   string
	s3Region   string
)

func RunBackup(ctx context.Context, job *scheduler.JobConfig) error {
	config := db.Config{
		Host:     job.DBHost,
		Port:     job.DBPort,
		User:     job.DBUser,
		Password: job.DBPassword,
		Database: job.DBName,
	}

	var database db.Database
	switch job.DBType {
	case "postgres":
		database = db.NewPostgresDB()
	default:
		return fmt.Errorf("unsupported database type: %s", job.DBType)
	}

	// Configure storage provider
	var store storage.Storage
	switch job.Storage {
	case "local":
		store = storage.NewLocalStorage(".")
	case "s3":
		var err error
		store, err = storage.NewS3Storage(ctx, job.S3Bucket, job.S3Region)
		if err != nil {
			return fmt.Errorf("error initializing S3 storage: %w", err)
		}
	default:
		return fmt.Errorf("unsupported storage destination: %s", job.Storage)
	}

	// Adjust output path if compression is enabled and missing .gz extension
	finalOutputPath := job.OutputPath
	if job.Compress && !strings.HasSuffix(finalOutputPath, ".gz") {
		finalOutputPath += ".gz"
	}

	// Set up pipeline using io.Pipe
	pr, pw := io.Pipe()

	// Run backup generation in a separate goroutine
	go func() {
		var err error
		var outWriter io.WriteCloser = pw

		if job.Compress {
			gzipProc := processor.NewGzipProcessor(pw)
			outWriter = gzipProc
		}

		err = database.Backup(ctx, config, outWriter)
		
		if job.Compress {
			outWriter.Close()
		}
		pw.CloseWithError(err)
	}()

	// The main thread reads from the pipe and saves to the storage
	err := store.Save(ctx, pr, finalOutputPath)
	if err != nil {
		return fmt.Errorf("backup failed during save: %w", err)
	}

	return nil
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a specific database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fmt.Printf("Starting backup for DB Type: %s, Name: %s\n", dbType, dbName)

		job := &scheduler.JobConfig{
			DBType:     dbType,
			DBName:     dbName,
			DBHost:     dbHost,
			DBPort:     dbPort,
			DBUser:     dbUser,
			DBPassword: dbPassword,
			Storage:    storageDst,
			OutputPath: outputPath,
			S3Bucket:   s3Bucket,
			S3Region:   s3Region,
			Compress:   compress,
		}

		err := RunBackup(ctx, job)
		if err != nil {
			fmt.Printf("Backup error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Backup completed successfully! Saved to [%s]\n", storageDst)
	},
}

func init() {
	// Attach flags to the backup command
	backupCmd.Flags().StringVarP(&dbType, "db", "d", "", "Database type (mysql, postgres, etc.)")
	backupCmd.Flags().StringVarP(&dbName, "name", "n", "", "Name of the database")
	backupCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	backupCmd.Flags().IntVar(&dbPort, "port", 0, "Database port (default depends on DB)")
	backupCmd.Flags().StringVar(&dbUser, "user", "", "Database user")
	backupCmd.Flags().StringVar(&dbPassword, "password", "", "Database password")
	backupCmd.Flags().StringVarP(&outputPath, "output", "o", "backup.sql", "Output path or key for the backup file")
	backupCmd.Flags().BoolVarP(&compress, "compress", "c", true, "Enable gzip compression")
	
	backupCmd.Flags().StringVar(&storageDst, "storage", "local", "Storage destination (local, s3)")
	backupCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (required if storage is s3)")
	backupCmd.Flags().StringVar(&s3Region, "s3-region", "", "AWS region for S3")

	// Mark flags as required
	backupCmd.MarkFlagRequired("db")
	backupCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(backupCmd)
}

