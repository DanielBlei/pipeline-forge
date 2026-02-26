/*
Copyright 2025 Daniel Blei.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	ResultInitialized      = "initialized"
	ResultValidationFailed = "validation_failed"
	ResultRequeued         = "requeued"
)

// StagingReconcileTotal counts Staging reconciliations partitioned by result.
var StagingReconcileTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "pipeline_forge_reconciliations_total",
		Help: "Total number of Staging reconciliations partitioned by result.",
	},
	[]string{"result"},
)

// Register registers all custom metrics with the controller-runtime registry.
// Safe to call multiple times — returns nil if already registered.
func Register() error {
	if err := ctrlmetrics.Registry.Register(StagingReconcileTotal); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return err
		}
	}
	return nil
}
