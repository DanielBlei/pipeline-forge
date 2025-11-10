package components

import (
	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultLimits = corev1.ResourceList{
		corev1.ResourceCPU:    resourceMustParse("500m"),
		corev1.ResourceMemory: resourceMustParse("256Mi"),
	}
	defaultRequests = corev1.ResourceList{
		corev1.ResourceCPU:    resourceMustParse("250m"),
		corev1.ResourceMemory: resourceMustParse("128Mi"),
	}
)

// resourceMustParse parse resource quantities
func resourceMustParse(val string) resource.Quantity {
	q, _ := resource.ParseQuantity(val)
	return q
}

// ensureResourceRequirements ensures that the resource requirements are not nil.
func ensureResourceRequirements(
	resources *corev1.ResourceRequirements,
	defaultLimits corev1.ResourceList,
	defaultRequests corev1.ResourceList,
) *corev1.ResourceRequirements {
	if resources == nil {
		return &corev1.ResourceRequirements{
			Limits:   defaultLimits,
			Requests: defaultRequests,
		}
	}
	return resources
}

// CreateCronJob creates a Kubernetes CronJob definition with the given parameters.
func CreateCronJob(
	name string,
	namespace string,
	schedule string,
	image string,
	command []string,
	resources *corev1.ResourceRequirements,
) *batchv1.CronJob {
	resources = ensureResourceRequirements(resources, defaultLimits, defaultRequests)

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
func CreateJob(
	name string,
	namespace string,
	image string,
	command []string,
	resources *corev1.ResourceRequirements,
) *batchv1.Job {
	resources = ensureResourceRequirements(resources, defaultLimits, defaultRequests)

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
							Resources: *resources,
						},
					},
				},
			},
		},
	}
}

// CreateTrigger creates a Kubernetes Trigger definition with the given parameters.
func CreateTrigger(
	name string,
	namespace string,
	image string,
	triggerType corev1alpha1.TriggerType,
	resources *corev1.ResourceRequirements,
	schedule string,
	description string,
	owner string,
) *corev1alpha1.Trigger {
	return &corev1alpha1.Trigger{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1alpha1.TriggerSpec{
			Name:        name,
			Type:        triggerType,
			Description: description,
			Owner:       owner,
			Image:       image,
			Resources: &corev1.ResourceRequirements{
				Limits:   resources.Limits,
				Requests: resources.Requests,
			},
			Schedule: schedule,
		},
	}
}

// CreateServiceAccount creates a Kubernetes ServiceAccount definition with the given parameters.
func CreateServiceAccount(name string, namespace string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
