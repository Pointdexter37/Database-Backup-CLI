package cmd

import (
	"context"
	"fmt"
	"os"

	"dbbackup/internal/db"

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

		err := database.Backup(context.Background(), config, outputPath)
		if err != nil {
			fmt.Printf("Backup failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Backup completed successfully! Saved to: %s\n", outputPath)
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

	// Mark flags as required
	backupCmd.MarkFlagRequired("db")
	backupCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(backupCmd)
}