package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"dbbackup/internal/db"
	"dbbackup/internal/processor"
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

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a specific database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fmt.Printf("Starting backup for DB Type: %s, Name: %s\n", dbType, dbName)

		config := db.Config{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
		}

		var database db.Database
		switch dbType {
		case "postgres":
			database = db.NewPostgresDB()
		default:
			fmt.Printf("Error: Unsupported database type: %s\n", dbType)
			os.Exit(1)
		}

		// Configure storage provider
		var store storage.Storage
		switch storageDst {
		case "local":
			store = storage.NewLocalStorage(".")
		case "s3":
			var err error
			store, err = storage.NewS3Storage(ctx, s3Bucket, s3Region)
			if err != nil {
				fmt.Printf("Error initializing S3 storage: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Error: Unsupported storage destination: %s\n", storageDst)
			os.Exit(1)
		}

		// Adjust output path if compression is enabled and missing .gz extension
		finalOutputPath := outputPath
		if compress && !strings.HasSuffix(finalOutputPath, ".gz") {
			finalOutputPath += ".gz"
		}

		// Set up pipeline using io.Pipe
		pr, pw := io.Pipe()

		// Run backup generation in a separate goroutine
		go func() {
			var err error
			var outWriter io.WriteCloser = pw

			if compress {
				gzipProc := processor.NewGzipProcessor(pw)
				outWriter = gzipProc
			}

			err = database.Backup(ctx, config, outWriter)
			
			// Close writers in correct order
			if compress {
				outWriter.Close()
			}
			pw.CloseWithError(err)
		}()

		// The main thread reads from the pipe and saves to the storage
		err := store.Save(ctx, pr, finalOutputPath)
		if err != nil {
			fmt.Printf("Backup failed during save: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Backup completed successfully! Saved to [%s]: %s\n", storageDst, finalOutputPath)
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

