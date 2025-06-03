package e2e_test

import (
	"fmt"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	rtServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
)

var _ = Describe("LocalMavenRepository", func() {

	When("a new local maven repository is created", func() {
		It("should exist in Artifactory", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-local-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a repository resource in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.LocalMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.LocalMavenRepositorySpec{
					ForProvider: v1alpha1.LocalMavenRepositoryParameters{
						Description: ptr.To("Test Local Maven Repository"),
					},
					ResourceSpec: v1.ResourceSpec{
						ProviderConfigReference: &v1.Reference{
							Name: "my-artifactory-providerconfig-write",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func(ctx SpecContext) {
				By("Deleting the repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.LocalMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.LocalMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())

				By("Verifying repository does not exist in Artifactory")
				repoDetails := rtServices.RepositoryDetails{}
				err = rtWriteClient.GetRepository(repoName, &repoDetails)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("400"))
				Expect(err.Error()).To(ContainSubstring("Bad Request"))
			})

			By("Waiting for the repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.LocalMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtWriteClient.GetRepository(repoName, &repoDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(repoDetails.Key).To(Equal(repoName))
			Expect(repoDetails.Description).To(Equal("Test Local Maven Repository"))
			Expect(repoDetails.GetRepoType()).To(Equal("local"))
			Expect(repoDetails.PackageType).To(Equal("maven"))
		})
	})
})
