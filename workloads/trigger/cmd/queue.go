package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/gcp" // Register GCP provider
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

var (
	topicName   string
	queueRegion string
)

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Trigger processing for message queue events",
	Long: `Monitor and process events from message queues across cloud providers.

This command triggers processing for messages published to the specified queue/topic.
Currently supports GCP Pub/Sub, with AWS SQS and Azure Service Bus coming soon.`,
	RunE: runQueue,
}

func init() {
	queueCmd.Flags().StringVarP(&topicName, "topic", "t", "", "Topic/queue name (required)")
	queueCmd.Flags().StringVarP(&queueRegion, "region", "r", "", "Topic/queue location/region")

	if err := queueCmd.MarkFlagRequired("topic"); err != nil {
		panic(err)
	}
	if err := queueCmd.MarkFlagRequired("region"); err != nil {
		panic(err)
	}
}

func runQueue(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	p, err := provider.GetQueueProvider(provider.Provider(cloudProvider))
	if err != nil {
		logger.Error("failed to get queue provider", zap.Error(err))
		return err
	}

	config := provider.QueueConfig{
		Provider:  provider.Provider(cloudProvider),
		ProjectID: projectID,
		Topic:     topicName,
		Region:    queueRegion,
		Verbose:   verbose,
	}

	return p.TriggerByQueue(ctx, logger, config)
}
