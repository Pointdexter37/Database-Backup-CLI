package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"dbbackup/internal/db"
	"dbbackup/internal/processor"

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
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a specific database",
	Run: func(cmd *cobra.Command, args []string) {
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

		// Adjust output path if compression is enabled and missing .gz extension
		finalOutputPath := outputPath
		if compress && !strings.HasSuffix(finalOutputPath, ".gz") {
			finalOutputPath += ".gz"
		}

		// Open output file
		file, err := os.Create(finalOutputPath)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		// Set up pipeline
		var outWriter io.Writer = file
		if compress {
			gzipProc := processor.NewGzipProcessor(file)
			defer gzipProc.Close()
			outWriter = gzipProc
		}

		err = database.Backup(context.Background(), config, outWriter)
		if err != nil {
			fmt.Printf("Backup failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Backup completed successfully! Saved to: %s\n", finalOutputPath)
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
	backupCmd.Flags().StringVarP(&outputPath, "output", "o", "backup.sql", "Output path for the backup file")
	backupCmd.Flags().BoolVarP(&compress, "compress", "c", true, "Enable gzip compression")

	// Mark flags as required
	backupCmd.MarkFlagRequired("db")
	backupCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(backupCmd)
}
