package e2e_test

import (
	"fmt"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
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

var _ = Describe("Assume control of an existing resource", func() {
	When("a local NPM repository already exists in Artifactory", func() {
		It("should assume control of that repository", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-local-npm-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a local repository resource in the Artifactory write instance")
			params := services.NewNpmLocalRepositoryParams()
			params.Key = repoName
			params.Description = "Test Local Npm Repository"
			err := rtWriteClient.CreateLocalRepository().Npm(params)
			Expect(err).NotTo(HaveOccurred())

			// Clean it up, if it still exists at the end of the test
			DeferCleanup(func(ctx SpecContext) {
				rtWriteClient.DeleteRepository(repoName)
			})

			By("Verifying the repository exists in Artifactory write instances")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtWriteClient.GetRepository(repoName, &repoDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(repoDetails.Key).To(Equal(repoName))
			Expect(repoDetails.Description).To(Equal("Test Local Npm Repository"))
			Expect(repoDetails.GetRepoType()).To(Equal("local"))
			Expect(repoDetails.PackageType).To(Equal("npm"))

			By("Creating a repository resource in Kubernetes")
			err = k8sClient.Create(ctx, &v1alpha1.LocalNpmRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.LocalNpmRepositorySpec{
					ForProvider: v1alpha1.LocalNpmRepositoryParameters{
						Description: ptr.To("Test Local Npm Repository"),
					},
					ResourceSpec: v1.ResourceSpec{
						ProviderConfigReference: &v1.Reference{
							Name: "my-artifactory-providerconfig-write",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for the repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.LocalNpmRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			By("Deleting the repository resource from Kubernetes")
			err = k8sClient.Delete(ctx, &v1alpha1.LocalNpmRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for the repository resource to be deleted")
			Eventually(func() bool {
				repo := &v1alpha1.LocalNpmRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				return errors.IsNotFound(err)
			}, "2m", "5s").Should(BeTrue())

			By("Ensuring that the repository has now been deleted from Artifactory")
			err = rtWriteClient.GetRepository(repoName, &repoDetails)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("400"))
			Expect(err.Error()).Should(ContainSubstring("Bad Request"))
		})
	})
})
