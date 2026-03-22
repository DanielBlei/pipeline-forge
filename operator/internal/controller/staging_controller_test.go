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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
			Namespace: "default",
		}
		cronJobNamespacedName := types.NamespacedName{
			Name:      cronJobName,
			Namespace: "default",
		}
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
			err = k8sClient.Get(ctx, typeNamespacedName, &corev1alpha1.Staging{})
			if err != nil && errors.IsNotFound(err) {
				resource := &corev1alpha1.Staging{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: corev1alpha1.StagingSpec{
						Description: "Testing Staging",
						Owner:       "Pipeline Forge",
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
			resource := &corev1alpha1.Staging{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("Cleanup the specific resource instance Staging")
				controllerutil.RemoveFinalizer(resource, stagingFinalizer)
				Expect(k8sClient.Update(ctx, resource)).To(Succeed())
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			By("Cleanup the CronJob")
			cronJob := &batchv1.CronJob{}
			err = k8sClient.Get(ctx, cronJobNamespacedName, cronJob)
			if err == nil {
				Expect(k8sClient.Delete(ctx, cronJob)).To(Succeed())
			}
		})

		It("should add finalizer on first reconcile", func() {
			controllerReconciler := &StagingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			By("Reconciling to add the finalizer")
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", 0), "should requeue after adding finalizer")

			By("Verifying the finalizer is present")
			updated := &corev1alpha1.Staging{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updated)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(updated, stagingFinalizer)).To(BeTrue())
		})

		It("should proceed with reconciliation after finalizer is set", func() {
			controllerReconciler := &StagingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			By("First reconcile adds finalizer")
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Second reconcile proceeds to domain logic")
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(Equal(defaultRequeueAfter),
				"should requeue with default interval after successful reconciliation")

			By("Verifying ingest status was updated")
			updated := &corev1alpha1.Staging{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updated)).To(Succeed())
			Expect(updated.Status.Ingest.Status).NotTo(BeNil())
			Expect(*updated.Status.Ingest.Status).To(Equal(corev1alpha1.ObjConditionReady))
		})

		It("should remove finalizer on deletion", func() {
			controllerReconciler := &StagingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			By("Reconciling to add the finalizer")
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Deleting the Staging object")
			toDelete := &corev1alpha1.Staging{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, toDelete)).To(Succeed())
			Expect(k8sClient.Delete(ctx, toDelete)).To(Succeed())

			By("Reconciling the deletion")
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero(), "should not requeue after deletion")

			By("Verifying the object is fully removed")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, &corev1alpha1.Staging{})
				return errors.IsNotFound(err)
			}).Should(BeTrue())
		})
	})

	Context("When reconciling with bootstrap mode", func() {
		const bootstrapName = "bootstrap-test"
		const bootstrapCronJobName = "bootstrap-cj"

		ctx := context.Background()

		bootstrapNamespacedName := types.NamespacedName{
			Name:      bootstrapName,
			Namespace: "default",
		}
		bootstrapCronJobNamespacedName := types.NamespacedName{
			Name:      bootstrapCronJobName,
			Namespace: "default",
		}

		BeforeEach(func() {
			By("creating the Staging resource with bootstrap mode")
			resource := &corev1alpha1.Staging{
				ObjectMeta: metav1.ObjectMeta{
					Name:      bootstrapName,
					Namespace: "default",
				},
				Spec: corev1alpha1.StagingSpec{
					Ingest: corev1alpha1.IngestSpec{
						Mode:     "bootstrap",
						Type:     "cronjob",
						Name:     bootstrapCronJobName,
						Schedule: "*/5 * * * *",
						Image:    "busybox:latest",
						Command:  []string{"echo", "ingest"},
					},
					Transform: corev1alpha1.TransformSpec{
						Project: "p",
						Target:  "dev",
						Image:   "dbt:latest",
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &corev1alpha1.Staging{}
			err := k8sClient.Get(ctx, bootstrapNamespacedName, resource)
			if err == nil {
				controllerutil.RemoveFinalizer(resource, stagingFinalizer)
				Expect(k8sClient.Update(ctx, resource)).To(Succeed())
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			cronJob := &batchv1.CronJob{}
			err = k8sClient.Get(ctx, bootstrapCronJobNamespacedName, cronJob)
			if err == nil {
				Expect(k8sClient.Delete(ctx, cronJob)).To(Succeed())
			}
		})

		It("should set owner reference on bootstrapped CronJob", func() {
			controllerReconciler := &StagingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			By("First reconcile adds finalizer")
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: bootstrapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Second reconcile creates CronJob with owner reference")
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: bootstrapNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the CronJob exists with correct owner reference")
			cronJob := &batchv1.CronJob{}
			Expect(k8sClient.Get(ctx, bootstrapCronJobNamespacedName, cronJob)).To(Succeed())
			Expect(cronJob.OwnerReferences).To(HaveLen(1))
			Expect(cronJob.OwnerReferences[0].Name).To(Equal(bootstrapName))
			Expect(cronJob.OwnerReferences[0].Kind).To(Equal("Staging"))
			Expect(*cronJob.OwnerReferences[0].Controller).To(BeTrue())
		})
	})
})
