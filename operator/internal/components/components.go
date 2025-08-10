package components

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// resourceMustParse parse resource quantities
func resourceMustParse(val string) resource.Quantity {
	q, _ := resource.ParseQuantity(val)
	return q
}

// CreateCronJob creates a Kubernetes CronJob definition with the given parameters.
func CreateCronJob(
	name string,
	namespace string,
	schedule string,
	image string,
	command []string,
	args []string,
	resources *corev1.ResourceRequirements,
) *batchv1.CronJob {
	if resources == nil {
		defaultResources := corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resourceMustParse("500m"),
				corev1.ResourceMemory: resourceMustParse("256Mi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resourceMustParse("250m"),
				corev1.ResourceMemory: resourceMustParse("128Mi"),
			},
		}
		resources = &defaultResources
	}

	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:      name,
									Image:     image,
									Command:   command,
									Args:      args,
									Resources: *resources,
								},
							},
						},
					},
				},
			},
		},
	}
}

// NewJob creates a Kubernetes Job definition with the given parameters.
func NewJob(
	name string,
	namespace string,
	image string,
	command []string,
	args []string,
	resources *corev1.ResourceRequirements,
) *batchv1.Job {
	if resources == nil {
		defaultResources := corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resourceMustParse("500m"),
				corev1.ResourceMemory: resourceMustParse("256Mi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resourceMustParse("250m"),
				corev1.ResourceMemory: resourceMustParse("128Mi"),
			},
		}
		resources = &defaultResources
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:      name,
							Image:     image,
							Command:   command,
							Args:      args,
							Resources: *resources,
						},
					},
				},
			},
		},
	}
}
