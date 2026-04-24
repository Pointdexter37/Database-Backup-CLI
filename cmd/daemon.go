package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"dbbackup/internal/logger"
	"dbbackup/internal/scheduler"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the background scheduling daemon",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		if err := logger.InitLogger("backup_daemon.log"); err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		logger.Log.Info("Starting DBBackup Daemon")

		manager := scheduler.NewManager("schedules.json")
		jobs, err := manager.LoadJobs()
		if err != nil {
			logger.Log.Error("Failed to load schedules", "error", err)
			os.Exit(1)
		}

		if len(jobs) == 0 {
			logger.Log.Info("No active schedules found. Exiting.")
			return
		}

		c := cron.New()

		for _, j := range jobs {
			jobCfg := j // capture loop variable
			_, err := c.AddFunc(jobCfg.CronExpr, func() {
				logger.Log.Info("Executing scheduled backup", "jobID", jobCfg.ID, "db", jobCfg.DBName)
				err := RunBackup(context.Background(), &jobCfg)
				if err != nil {
					logger.Log.Error("Scheduled backup failed", "jobID", jobCfg.ID, "error", err)
				} else {
					logger.Log.Info("Scheduled backup succeeded", "jobID", jobCfg.ID)
				}
			})
			if err != nil {
				logger.Log.Error("Failed to schedule job", "jobID", jobCfg.ID, "error", err)
			} else {
				logger.Log.Info("Scheduled job", "jobID", jobCfg.ID, "cron", jobCfg.CronExpr)
			}
		}

		c.Start()

		// Wait for termination signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Log.Info("Shutting down daemon...")
		c.Stop()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
