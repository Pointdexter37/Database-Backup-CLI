package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	dbType string
	dbName string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a specific database",
	Run: func(cmd *cobra.Command, args []string) {
		// This is where we will later call internal/db logic
		fmt.Printf("Starting backup for DB Type: %s, Name: %s\n", dbType, dbName)
	},
}

func init() {
	// Attach flags to the backup command
	backupCmd.Flags().StringVarP(&dbType, "db", "d", "", "Database type (mysql, postgres, etc.)")
	backupCmd.Flags().StringVarP(&dbName, "name", "n", "", "Name of the database")
	
	// Mark flags as required
	backupCmd.MarkFlagRequired("db")
	backupCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(backupCmd)
}