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

	"github.com/myorg/provider-jfrogartifactory/apis/repository/v1alpha1"
)

var _ = Describe("E2E Tests", func() {

	Describe("GenericRepository", func() {
		When("a new repository is created", func() {
			It("should exist in Artifactory", func(ctx SpecContext) {
				By("Creating a repository resource in Kubernetes")
				err := k8sClient.Create(ctx, &v1alpha1.GenericRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-repo",
					},
					Spec: v1alpha1.GenericRepositorySpec{
						ForProvider: v1alpha1.GenericRepositoryParameters{
							Description: ptr.To("Test repository"),
						},
						ResourceSpec: v1.ResourceSpec{
							ProviderConfigReference: &v1.Reference{
								Name: "my-artifactory-providerconfig-read",
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				DeferCleanup(func(ctx SpecContext) {
					By("Deleting the repository resource from Kubernetes")
					err := k8sClient.Delete(ctx, &v1alpha1.GenericRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-repo",
						},
					})
					Expect(err).NotTo(HaveOccurred())

					By("Waiting for the repository resource to be deleted")
					Eventually(func() bool {
						repo := &v1alpha1.GenericRepository{}
						err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-repo"}, repo)
						return errors.IsNotFound(err)
					}, "2m", "5s").Should(BeTrue())
				})

				By("Waiting for the repository to be ready in Kubernetes")
				Eventually(func() bool {
					repo := &v1alpha1.GenericRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-repo"}, repo)
					Expect(err).NotTo(HaveOccurred())
					return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
						repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
				}, "2m", "5s").Should(BeTrue())

				By("Verifying the repository exists in Artifactory")
				repoDetails := rtServices.RepositoryDetails{}
				err = rtReadClient.GetRepository("test-repo", &repoDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(repoDetails.Key).To(Equal("test-repo"))
				Expect(repoDetails.Description).To(Equal("Test repository"))
				Expect(repoDetails.GetRepoType()).To(Equal("local"))
				Expect(repoDetails.PackageType).To(Equal("generic"))
			})
		})
	})
})
