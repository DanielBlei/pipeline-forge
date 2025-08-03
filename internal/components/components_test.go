package components

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func resourceRequirementsEqual(a, b corev1.ResourceRequirements) bool {
	return a.Limits.Cpu().Cmp(*b.Limits.Cpu()) == 0 &&
		a.Limits.Memory().Cmp(*b.Limits.Memory()) == 0 &&
		a.Requests.Cpu().Cmp(*b.Requests.Cpu()) == 0 &&
		a.Requests.Memory().Cmp(*b.Requests.Memory()) == 0
}

func TestCreateCronJob(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			name      string
			namespace string
			schedule  string
			image     string
			command   []string
			args      []string
			resources *corev1.ResourceRequirements
		}
		wantContainer corev1.Container
		wantSchedule  string
	}{
		{
			name: "WithResources",
			args: struct {
				name      string
				namespace string
				schedule  string
				image     string
				command   []string
				args      []string
				resources *corev1.ResourceRequirements
			}{
				name:      "test-cronjob",
				namespace: "test-ns",
				schedule:  "*/5 * * * *",
				image:     "busybox",
				command:   []string{"echo"},
				args:      []string{"hello"},
				resources: &corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
			wantContainer: corev1.Container{
				Name:    "test-cronjob",
				Image:   "busybox",
				Command: []string{"echo"},
				Args:    []string{"hello"},
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
			wantSchedule: "*/5 * * * *",
		},
		{
			name: "DefaultResources",
			args: struct {
				name      string
				namespace string
				schedule  string
				image     string
				command   []string
				args      []string
				resources *corev1.ResourceRequirements
			}{
				name:      "test-default",
				namespace: "default",
				schedule:  "0 0 * * *",
				image:     "alpine",
				command:   nil,
				args:      nil,
				resources: nil,
			},
			wantContainer: corev1.Container{
				Name:  "test-default",
				Image: "alpine",
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("250m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
			wantSchedule: "0 0 * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cronJob := CreateCronJob(
				tt.args.name,
				tt.args.namespace,
				tt.args.schedule,
				tt.args.image,
				tt.args.command,
				tt.args.args,
				tt.args.resources,
			)
			if cronJob.ObjectMeta.Name != tt.args.name {
				t.Errorf("expected name %q, got %q", tt.args.name, cronJob.ObjectMeta.Name)
			}
			if cronJob.ObjectMeta.Namespace != tt.args.namespace {
				t.Errorf("expected namespace %q, got %q", tt.args.namespace, cronJob.ObjectMeta.Namespace)
			}
			if cronJob.Spec.Schedule != tt.wantSchedule {
				t.Errorf("expected schedule %q, got %q", tt.wantSchedule, cronJob.Spec.Schedule)
			}
			container := cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
			if container.Image != tt.wantContainer.Image {
				t.Errorf("expected image %q, got %q", tt.wantContainer.Image, container.Image)
			}
			if len(tt.wantContainer.Command) > 0 && !equalStringSlice(container.Command, tt.wantContainer.Command) {
				t.Errorf("expected command %v, got %v", tt.wantContainer.Command, container.Command)
			}
			if len(tt.wantContainer.Args) > 0 && !equalStringSlice(container.Args, tt.wantContainer.Args) {
				t.Errorf("expected args %v, got %v", tt.wantContainer.Args, container.Args)
			}
			if !resourceRequirementsEqual(container.Resources, tt.wantContainer.Resources) {
				t.Errorf("expected resources %+v, got %+v", tt.wantContainer.Resources, container.Resources)
			}
			if cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyOnFailure {
				t.Errorf("expected restart policy OnFailure, got %s", cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy)
			}
			if container.Name != tt.args.name {
				t.Errorf("expected container name %q, got %q", tt.args.name, container.Name)
			}
		})
	}
}

func TestNewJob(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			name      string
			namespace string
			image     string
			command   []string
			args      []string
			resources *corev1.ResourceRequirements
		}
		wantContainer corev1.Container
	}{
		{
			name: "WithResources",
			args: struct {
				name      string
				namespace string
				image     string
				command   []string
				args      []string
				resources *corev1.ResourceRequirements
			}{
				name:      "test-job",
				namespace: "ns1",
				image:     "ubuntu",
				command:   []string{"/bin/sh"},
				args:      []string{"-c", "echo world"},
				resources: &corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("2"),
						corev1.ResourceMemory: resource.MustParse("1Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				},
			},
			wantContainer: corev1.Container{
				Name:    "test-job",
				Image:   "ubuntu",
				Command: []string{"/bin/sh"},
				Args:    []string{"-c", "echo world"},
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("2"),
						corev1.ResourceMemory: resource.MustParse("1Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				},
			},
		},
		{
			name: "DefaultResources",
			args: struct {
				name      string
				namespace string
				image     string
				command   []string
				args      []string
				resources *corev1.ResourceRequirements
			}{
				name:      "job-default",
				namespace: "ns2",
				image:     "alpine",
				command:   nil,
				args:      nil,
				resources: nil,
			},
			wantContainer: corev1.Container{
				Name:  "job-default",
				Image: "alpine",
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("250m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := NewJob(
				tt.args.name,
				tt.args.namespace,
				tt.args.image,
				tt.args.command,
				tt.args.args,
				tt.args.resources,
			)
			if job.ObjectMeta.Name != tt.args.name {
				t.Errorf("expected name %q, got %q", tt.args.name, job.ObjectMeta.Name)
			}
			if job.ObjectMeta.Namespace != tt.args.namespace {
				t.Errorf("expected namespace %q, got %q", tt.args.namespace, job.ObjectMeta.Namespace)
			}
			container := job.Spec.Template.Spec.Containers[0]
			if container.Image != tt.wantContainer.Image {
				t.Errorf("expected image %q, got %q", tt.wantContainer.Image, container.Image)
			}
			if len(tt.wantContainer.Command) > 0 && !equalStringSlice(container.Command, tt.wantContainer.Command) {
				t.Errorf("expected command %v, got %v", tt.wantContainer.Command, container.Command)
			}
			if len(tt.wantContainer.Args) > 0 && !equalStringSlice(container.Args, tt.wantContainer.Args) {
				t.Errorf("expected args %v, got %v", tt.wantContainer.Args, container.Args)
			}
			if !resourceRequirementsEqual(container.Resources, tt.wantContainer.Resources) {
				t.Errorf("expected resources %+v, got %+v", tt.wantContainer.Resources, container.Resources)
			}
			if job.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyOnFailure {
				t.Errorf("expected restart policy OnFailure, got %s", job.Spec.Template.Spec.RestartPolicy)
			}
			if container.Name != tt.args.name {
				t.Errorf("expected container name %q, got %q", tt.args.name, container.Name)
			}
		})
	}
}

// Helper for comparing string slices
func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
