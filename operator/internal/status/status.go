package status

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateStatus updates only status that differ from objDeepCopy
func UpdateStatus(
	ctx context.Context,
	c client.Client,
	obj client.Object,
	objDeepCopy client.Object,
) error {
	patch := client.MergeFrom(objDeepCopy)
	return c.Status().Patch(ctx, obj, patch)
}
