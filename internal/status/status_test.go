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

package status

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	v1alpha1 "github.com/DanielBlei/pipeline-forge/api/v1alpha1"
)

var _ = Describe("Status Update", func() {
	Context("When updating status", func() {
		const resourceName = "test-update-status"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}

		BeforeEach(func() {
			By("creating the custom resource for testing status updates")
			resource := &v1alpha1.Trigger{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: v1alpha1.TriggerSpec{
					Type: "bigquery",
					BigQuery: &v1alpha1.BigQueryTriggerSpec{
						Project: "gcp-project",
						Dataset: "dataset_id",
						Table:   "table_id",
					},
				},
			}
			Expect(k8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			By("Cleanup the test resource")
			resource := &v1alpha1.Trigger{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully update the status", func() {
			By("Getting the created resource")
			resource := &v1alpha1.Trigger{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Creating a deep copy for the patch")
			objDeepCopy := resource.DeepCopy()

			By("Updating the status")
			resource.Status = v1alpha1.StatusCompleted

			By("Calling UpdateStatus function")
			err = UpdateStatus(ctx, k8sClient, resource, objDeepCopy)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the status was updated")
			updatedObj := &v1alpha1.Trigger{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedObj)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedObj.Status).To(Equal(v1alpha1.StatusCompleted))
		})
	})
})
