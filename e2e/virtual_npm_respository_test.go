package e2e_test

import (
	"fmt"
	"strings"

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

var _ = Describe("VirtualNpmRepository", func() {
	var localRepoName string

	BeforeEach(func(ctx SpecContext) {
		localRepoName = fmt.Sprintf("test-local-npm-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		By("Creating a local repository resource with read ProviderConfig in Kubernetes")
		err := k8sClient.Create(ctx, &v1alpha1.LocalNpmRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: localRepoName,
			},
			Spec: v1alpha1.LocalNpmRepositorySpec{
				ForProvider: v1alpha1.LocalNpmRepositoryParameters{
					Description: ptr.To("Test Local Npm Read Repository"),
				},
				ResourceSpec: v1.ResourceSpec{
					ProviderConfigReference: &v1.Reference{
						Name: "my-artifactory-providerconfig-read",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the local repository to be ready in Kubernetes")
		Eventually(func() bool {
			repo := &v1alpha1.LocalNpmRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: localRepoName}, repo)
			Expect(err).NotTo(HaveOccurred())
			return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
				repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
		}, "2m", "5s").Should(BeTrue())
		// Check for the actual existence as well
		By("Verifying the repository exists in Artifactory")
		repoDetails := rtServices.RepositoryDetails{}
		err = rtReadClient.GetRepository(localRepoName, &repoDetails)
		Expect(err).NotTo(HaveOccurred())
		Expect(repoDetails.Key).To(Equal(localRepoName))
		Expect(repoDetails.Description).To(Equal("Test Local Npm Read Repository"))
		Expect(repoDetails.GetRepoType()).To(Equal("local"))
		Expect(repoDetails.PackageType).To(Equal("npm"))
	})

	AfterEach(func(ctx SpecContext) {
		By("Deleting the local repository resource from Kubernetes")
		err := k8sClient.Delete(ctx, &v1alpha1.LocalNpmRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: localRepoName,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the local repository resource to be deleted")
		Eventually(func() bool {
			repo := &v1alpha1.LocalNpmRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: localRepoName}, repo)
			return errors.IsNotFound(err)
		}, "2m", "5s").Should(BeTrue())
	})

	When("a new repository is created", func() {
		It("should exist in Artifactory read instance", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-virtual-npm-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a virtual repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.VirtualNpmRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.VirtualNpmRepositorySpec{
					ForProvider: v1alpha1.VirtualNpmRepositoryParameters{
						Description: ptr.To("Test Virtual Npm Read Repository"),
						Repositories: []*string{
							ptr.To(localRepoName),
						},
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
				By("Deleting the virtual repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.VirtualNpmRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the virtual repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.VirtualNpmRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the virtual repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.VirtualNpmRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the virtual repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository(repoName, &repoDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(repoDetails.Key).To(Equal(repoName))
			Expect(repoDetails.Description).To(Equal("Test Virtual Npm Read Repository"))
			Expect(repoDetails.GetRepoType()).To(Equal("virtual"))
			Expect(repoDetails.PackageType).To(Equal("npm"))
		})
	})

	When("a new repository is created with non existing local repo in read instance", func() {
		It("should not exist in Artifactory read instance", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-virtual-npm-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a virtual repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.VirtualNpmRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.VirtualNpmRepositorySpec{
					ForProvider: v1alpha1.VirtualNpmRepositoryParameters{
						Description: ptr.To("Test Virtual Npm Read Repository"),
						Repositories: []*string{
							ptr.To("test-local-npm-read-repo-nonexistent"),
						},
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
				By("Deleting the virtual repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.VirtualNpmRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the virtual repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.VirtualNpmRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the virtual repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.VirtualNpmRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionFalse &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionFalse &&
					strings.Contains(repo.Status.GetCondition(v1.TypeSynced).Message,
						"Repository test-local-npm-read-repo-nonexistent does not exist")
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the virtual repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository(repoName, &repoDetails)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("400"))
			Expect(err.Error()).To(ContainSubstring("Bad Request"))
		})
	})
})
