package e2e_test

import (
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	rtServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/myorg/provider-jfrogartifactory/apis/artifactorygroup/v1alpha1"
)

var _ = Describe("Artifactory Group", func() {

	BeforeAll()

	AfterAll()

	When("a new group is created", func() {
		It("should exists in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory Group
			groupName := ""
			By("Creating an artifactory group resource in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.ArtifactoryGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: groupName, 
				},
				Spec: v1alpha1.ArtifactoryGroupSpec{
					ForProvider: v1alpha1.ArtifactoryGroupParameters{
						Name: &groupName,
						Description: ptr.To("Test Artifactory User"),
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Check if object is ready
			By("Waiting for group to be ready in Kubernetes")
			Eventually(func() bool {
				group := &v1alpha1.ArtifactoryGroup{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: groupName}, group)
				Expect(err).NotTo(HaveOccurred())
				return group.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					group.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			// Verify if group is created on Artifactory
			By("Verifying group exists on Artifactory")
			groupDetails := rtServices.GroupParams{}
			retrievedGroup, err := rtReadClient.GetGroup(groupDetails)
			Expect(retrievedGroup.Name).To(Equal(groupName))
			Expect(retrievedGroup.Description).To(Equal(""))


		})
	})

	When("a new group is created with invalid credentials", func() {
		It("should not exist in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory Group
			By("Creating an artifactory group resource in Kubernetes")

			// Check if object is not ready
			By("Waiting for group to fail in Kubernetes")

			// Verify if group is not created on Artifactory
			By("Verifying group does not exist on Artifactory")
		})
	})
})