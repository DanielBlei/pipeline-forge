package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/gcp" // Register GCP provider
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

var (
	datasetName string
	tableName   string
	tableRegion string
)

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Trigger processing for database table events",
	Long: `Monitor and process changes to database tables across cloud providers.
	
This command triggers processing for events occurring in the specified database table.
Currently supports GCP BigQuery, with AWS Redshift and Azure Synapse coming soon.`,
	RunE: runTable,
}

func init() {
	tableCmd.Flags().StringVarP(&datasetName, "dataset", "d", "", "Dataset/database name (required)")
	tableCmd.Flags().StringVarP(&tableName, "table", "t", "", "Table name (required)")
	tableCmd.Flags().StringVarP(&tableRegion, "region", "r", "", "Dataset/database location/region")

	if err := tableCmd.MarkFlagRequired("dataset"); err != nil {
		panic(err)
	}
	if err := tableCmd.MarkFlagRequired("table"); err != nil {
		panic(err)
	}
	tableCmd.MarkFlagsRequiredTogether("dataset", "table")
}

func runTable(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	p, err := provider.GetTableProvider(provider.Provider(cloudProvider))
	if err != nil {
		logger.Error("failed to get table provider", zap.Error(err))
		return err
	}

	config := provider.TableConfig{
		Provider:  provider.Provider(cloudProvider),
		ProjectID: projectID,
		Dataset:   datasetName,
		Table:     tableName,
		Region:    tableRegion,
		Verbose:   verbose,
	}

	return p.TriggerByTable(ctx, logger, config)
}
