package cmd

import (
	"fmt"
	"os"
	"time"

	"dbbackup/internal/scheduler"

	"github.com/spf13/cobra"
)

var (
	cronExpr string
	jobID    string
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage backup schedules",
}

var scheduleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new backup schedule",
	Run: func(cmd *cobra.Command, args []string) {
		manager := scheduler.NewManager("schedules.json")
		job := scheduler.JobConfig{
			ID:              fmt.Sprintf("job-%d", time.Now().UnixNano()),
			CronExpr:        cronExpr,
			DBType:          dbType,
			DBName:          dbName,
			DBHost:          dbHost,
			DBPort:          dbPort,
			DBUser:          dbUser,
			DBPassword:      dbPassword,
			DockerContainer: dockerContainer,
			Storage:         storageDst,
			OutputPath:      outputPath,
			S3Bucket:        s3Bucket,
			S3Region:        s3Region,
			Compress:        compress,
		}

		if err := manager.AddJob(job); err != nil {
			fmt.Printf("Failed to add schedule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added schedule %s with cron expression '%s'\n", job.ID, job.CronExpr)
	},
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backup schedules",
	Run: func(cmd *cobra.Command, args []string) {
		manager := scheduler.NewManager("schedules.json")
		jobs, err := manager.LoadJobs()
		if err != nil {
			fmt.Printf("Failed to list schedules: %v\n", err)
			os.Exit(1)
		}

		if len(jobs) == 0 {
			fmt.Println("No active schedules found.")
			return
		}

		for _, j := range jobs {
			fmt.Printf("- ID: %s | Cron: %s | DB: %s (%s) | Storage: %s\n", j.ID, j.CronExpr, j.DBName, j.DBType, j.Storage)
		}
	},
}

var scheduleRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a backup schedule by ID",
	Run: func(cmd *cobra.Command, args []string) {
		manager := scheduler.NewManager("schedules.json")
		if err := manager.RemoveJob(jobID); err != nil {
			fmt.Printf("Failed to remove schedule: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully removed schedule %s\n", jobID)
	},
}

func init() {
	scheduleAddCmd.Flags().StringVar(&cronExpr, "cron", "", "Cron expression (e.g., '0 2 * * *')")
	scheduleAddCmd.MarkFlagRequired("cron")

	// Inherit database connection flags from backup.go vars for simplicity in this implementation
	scheduleAddCmd.Flags().StringVarP(&dbType, "db", "d", "", "Database type")
	scheduleAddCmd.Flags().StringVarP(&dbName, "name", "n", "", "Name of the database")
	scheduleAddCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	scheduleAddCmd.Flags().IntVar(&dbPort, "port", 0, "Database port")
	scheduleAddCmd.Flags().StringVar(&dbUser, "user", "", "Database user")
	scheduleAddCmd.Flags().StringVar(&dbPassword, "password", "", "Database password")
	scheduleAddCmd.Flags().StringVarP(&outputPath, "output", "o", "backup.sql", "Output path or key")
	scheduleAddCmd.Flags().BoolVarP(&compress, "compress", "c", true, "Enable compression")
	scheduleAddCmd.Flags().StringVar(&storageDst, "storage", "local", "Storage destination")
	scheduleAddCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name")
	scheduleAddCmd.Flags().StringVar(&s3Region, "s3-region", "", "AWS region")
	scheduleAddCmd.Flags().StringVar(&dockerContainer, "docker-container", "", "Execute backup inside a specific docker container")

	scheduleAddCmd.MarkFlagRequired("db")
	scheduleAddCmd.MarkFlagRequired("name")

	scheduleRemoveCmd.Flags().StringVar(&jobID, "id", "", "ID of the schedule to remove")
	scheduleRemoveCmd.MarkFlagRequired("id")

	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	scheduleCmd.AddCommand(scheduleRemoveCmd)

	rootCmd.AddCommand(scheduleCmd)
}
