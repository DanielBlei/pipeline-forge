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

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/api/v1alpha1"
	"github.com/DanielBlei/pipeline-forge/test/utils"
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
			resource := &corev1alpha1.Trigger{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: corev1alpha1.TriggerSpec{
					Kind: "bigquery",
					BigQuery: &corev1alpha1.BigQueryTriggerSpec{
						Project: "gcp-project",
						Dataset: "dataset_id",
						Table:   "table_id",
					},
				},
			}
			Expect(utils.K8sClient.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			By("Cleanup the test resource")
			resource := &corev1alpha1.Trigger{}
			err := utils.K8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(utils.K8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully update the status", func() {
			By("Getting the created resource")
			resource := &corev1alpha1.Trigger{}
			err := utils.K8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Creating a deep copy for the patch")
			objDeepCopy := resource.DeepCopy()

			By("Updating the status")
			resource.Status.Status = "Completed"

			By("Calling UpdateStatus function")
			err = UpdateStatus(ctx, utils.K8sClient, resource, objDeepCopy)
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the status was updated")
			updatedObj := &corev1alpha1.Trigger{}
			err = utils.K8sClient.Get(ctx, typeNamespacedName, updatedObj)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedObj.Status.Status).To(Equal("Completed"))
		})
	})
})
