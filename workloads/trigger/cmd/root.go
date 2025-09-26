package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/bucket"
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/queue"
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/warehouse"
)

var (
	logger     *zap.Logger
	verbose    bool
	projectID  string
	bucketName string
	queueName  string
	dataset    string
	table      string
)

// rootCmd is the root command of the trigger workload CLI
var rootCmd = &cobra.Command{
	Use:     "trigger-workload",
	Version: "0.1.0",
	Short:   "Trigger workload",
	Long:    "Trigger workload process events from GCS, Pub/Sub, or BigQuery",
	Run: func(cmd *cobra.Command, args []string) {
		routeToSubcommand()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

// init initializes the root command
func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "project ID")
	rootCmd.PersistentFlags().StringVarP(&bucketName, "bucket", "b", "", "bucket name")
	rootCmd.PersistentFlags().StringVarP(&queueName, "queue", "q", "", "queue name")
	rootCmd.PersistentFlags().StringVarP(&dataset, "dataset", "d", "", "dataset name")
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", "", "table name")

	// Set flags workflow
	rootCmd.MarkPersistentFlagRequired("project")
	rootCmd.MarkFlagsMutuallyExclusive("bucket", "queue", "dataset")
	rootCmd.MarkFlagsRequiredTogether("dataset", "table")
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

// routeToSubcommand routes to the appropriate trigger function based on flags
func routeToSubcommand() {
	// Determine which resource type is specified
	if bucketName != "" {
		logger.Debug("triggering by bucket", zap.String("bucket", bucketName))
		bucket.TriggerByBucket(logger, projectID, bucketName, verbose)
	} else if queueName != "" {
		logger.Debug("triggering by queue", zap.String("queue", queueName))
		queue.TriggerByQueue(logger, projectID, queueName, verbose)
	} else if dataset != "" {
		logger.Debug("triggering by warehouse", zap.String("dataset", dataset), zap.String("table", table))
		warehouse.TriggerByWarehouse(logger, projectID, dataset, table, verbose)
	} else {
		logger.Error("no resource specified, Cobra failed to mark the flags as required")
	}
}
