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

package v1alpha1

import "time"

// DefaultPollIntervalSeconds is the default poll interval when not specified in the spec.
const DefaultPollIntervalSeconds int32 = 120

// IngestNamespace returns the namespace of the ingest resource,
// falling back to the Staging object's own namespace if not specified.
func (s *Staging) IngestNamespace() string {
	if s.Spec.Ingest.Namespace != "" {
		return s.Spec.Ingest.Namespace
	}
	return s.Namespace
}

// IngestPollInterval returns the poll interval as a time.Duration,
// defaulting to 120s if not specified.
func (s *Staging) IngestPollInterval() time.Duration {
	if s.Spec.Ingest.PollIntervalSeconds > 0 {
		return time.Duration(s.Spec.Ingest.PollIntervalSeconds) * time.Second
	}
	return time.Duration(DefaultPollIntervalSeconds) * time.Second
}

// IngestSuspended returns whether the ingest step is suspended.
func (s *Staging) IngestSuspended() bool {
	return s.Spec.Ingest.Suspend
}

// TransformSuspended returns whether the transform step is suspended.
func (s *Staging) TransformSuspended() bool {
	return s.Spec.Transform.Suspend
}
