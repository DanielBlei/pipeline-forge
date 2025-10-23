package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/gcp" // Register GCP provider
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

var (
	bucketName   string
	bucketRegion string
)

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Trigger processing for object storage bucket events",
	Long: `Monitor and process events from object storage buckets across cloud providers.

This command triggers processing for events occurring in the specified object storage bucket.,
Currently supports GCP Cloud Storage, with AWS S3 and Azure Blob Storage coming soon.`,
	RunE: runBucket,
}

func init() {
	bucketCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "Bucket name (required)")
	bucketCmd.Flags().StringVarP(&bucketRegion, "region", "r", "", "Bucket location/region")

	if err := bucketCmd.MarkFlagRequired("bucket"); err != nil {
		panic(err)
	}
	if err := bucketCmd.MarkFlagRequired("region"); err != nil {
		panic(err)
	}
}

func runBucket(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	p, err := provider.GetStorageProvider(provider.Provider(cloudProvider))
	if err != nil {
		logger.Error("failed to get storage provider", zap.Error(err))
		return err
	}

	config := provider.StorageConfig{
		Provider:  provider.Provider(cloudProvider),
		ProjectID: projectID,
		Bucket:    bucketName,
		Region:    bucketRegion,
		Verbose:   verbose,
	}

	return p.TriggerByBucket(ctx, logger, config)
}
