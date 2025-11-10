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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	"github.com/DanielBlei/pipeline-forge/operator/internal/components"
)

var _ = Describe("Staging Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"
		const cronJobName = "example-cronjob"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		cronJobNamespacedName := types.NamespacedName{
			Name:      cronJobName,
			Namespace: "default",
		}
		staging := &corev1alpha1.Staging{}

		BeforeEach(func() {
			By("creating the referenced CronJob")
			cronJob := &batchv1.CronJob{}
			err := k8sClient.Get(ctx, cronJobNamespacedName, cronJob)
			if err != nil && errors.IsNotFound(err) {
				cronJob = components.CreateCronJob(
					cronJobName,
					"default",
					"*/5 * * * *",
					"busybox:latest",
					[]string{"echo", "test"},
					nil,
				)
				Expect(k8sClient.Create(ctx, cronJob)).To(Succeed())
			}

			By("creating the custom resource for the Kind Staging")
			err = k8sClient.Get(ctx, typeNamespacedName, staging)
			if err != nil && errors.IsNotFound(err) {
				resource := &corev1alpha1.Staging{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: corev1alpha1.StagingSpec{
						Description: "Testing Staging",
						Owner:       "Daniel Blei",
						Ingest: corev1alpha1.IngestSpec{
							Mode:      "reference",
							Type:      "cronjob",
							Name:      cronJobName,
							Namespace: "default",
						},
						Transform: corev1alpha1.TransformSpec{
							Project: "example-dbt-project",
							Target:  "dev",
							Image:   "ghcr.io/example/dbt:latest",
							Models:  []string{"stg_model"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &corev1alpha1.Staging{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Staging")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())

			By("Cleanup the CronJob")
			cronJob := &batchv1.CronJob{}
			err = k8sClient.Get(ctx, cronJobNamespacedName, cronJob)
			if err == nil {
				Expect(k8sClient.Delete(ctx, cronJob)).To(Succeed())
			}
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &StagingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
