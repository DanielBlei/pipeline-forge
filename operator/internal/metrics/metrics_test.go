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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRegister_Idempotent(t *testing.T) {
	if err := Register(); err != nil {
		t.Fatalf("first Register() returned error: %v", err)
	}
	if err := Register(); err != nil {
		t.Fatalf("second Register() returned error: %v", err)
	}
}

func TestStagingReconcileTotal_EachLabelIncrements(t *testing.T) {
	// Use a fresh registry so test state is isolated from the global controller-runtime registry.
	reg := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pipeline_forge_staging_reconciliations_total_test",
			Help: "Test counter.",
		},
		[]string{"result"},
	)
	reg.MustRegister(counter)

	results := []string{ResultSkipped, ResultInitialized, ResultValidationFailed, ResultRequeued}
	for _, result := range results {
		t.Run(result, func(t *testing.T) {
			counter.WithLabelValues(result).Inc()
			if got := testutil.ToFloat64(counter.WithLabelValues(result)); got != 1 {
				t.Errorf("result=%q: expected 1, got %v", result, got)
			}
		})
	}
}
