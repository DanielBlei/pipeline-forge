package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logger        *zap.Logger
	verbose       bool
	projectID     string
	cloudProvider string
)

// rootCmd is the root command of the trigger workload CLI
var rootCmd = &cobra.Command{
	Use:     "trigger-workload",
	Version: "0.1.0",
	Short:   "Trigger workload for multi-cloud services",
	Long: `Trigger workload processes events from object storage, 
	message queues, and database tables across cloud providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help when no subcommand is provided
		if err := cmd.Help(); err != nil {
			panic(fmt.Sprintf("failed to display help: %v", err))
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
	// Hide completion command from help but keep it functional
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// init initializes the root command
func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&cloudProvider, "cloud-provider", "c", "", "Cloud provider (required)")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "Project ID (required)")

	// Mark flags as required
	if err := rootCmd.MarkPersistentFlagRequired("cloud-provider"); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("project"); err != nil {
		panic(err)
	}

	// Add subcommands
	rootCmd.AddCommand(bucketCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(tableCmd)
}

// initLogger initializes the logger
func initLogger() {
	var err error
	if verbose {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
}

// Execute executes the root command
func Execute() {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic(fmt.Sprintf("failed to initialize logger: %v", err))
		}
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Error("command execution failed", zap.Error(err))
		os.Exit(1)
	}
}
