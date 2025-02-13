package e2e_test

import (
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

var _ = Describe("VirtualMavenRepository", Ordered, func() {

	BeforeAll(func(ctx SpecContext) {
		By("Creating a local repository resource with read ProviderConfig in Kubernetes")
		err := k8sClient.Create(ctx, &v1alpha1.LocalMavenRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-local-maven-read-repo",
			},
			Spec: v1alpha1.LocalMavenRepositorySpec{
				ForProvider: v1alpha1.LocalMavenRepositoryParameters{
					Description: ptr.To("Test Local Maven Read Repository"),
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
			repo := &v1alpha1.LocalMavenRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-local-maven-read-repo"}, repo)
			Expect(err).NotTo(HaveOccurred())
			return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
				repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
		}, "2m", "5s").Should(BeTrue())
	})

	AfterAll(func(ctx SpecContext) {
		By("Deleting the local repository resource from Kubernetes")
		err := k8sClient.Delete(ctx, &v1alpha1.LocalMavenRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-local-maven-read-repo",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the local repository resource to be deleted")
		Eventually(func() bool {
			repo := &v1alpha1.LocalMavenRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-local-maven-read-repo"}, repo)
			return errors.IsNotFound(err)
		}, "2m", "5s").Should(BeTrue())
	})

	When("a new repository is created", func() {
		It("should exist in Artifactory read instance", func(ctx SpecContext) {
			By("Creating a virtual repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.VirtualMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-virtual-maven-repo-read",
				},
				Spec: v1alpha1.VirtualMavenRepositorySpec{
					ForProvider: v1alpha1.VirtualMavenRepositoryParameters{
						Description:   ptr.To("Test Virtual Maven Read Repository"),
						RepoLayoutRef: ptr.To("maven-2-default"),
						Repositories: []*string{
							ptr.To("test-local-maven-read-repo"),
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
				err := k8sClient.Delete(ctx, &v1alpha1.VirtualMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-virtual-maven-repo-read",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the virtual repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.VirtualMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-virtual-maven-repo-read"}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the virtual repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.VirtualMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-virtual-maven-repo-read"}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the virtual repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository("test-virtual-maven-repo-read", &repoDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(repoDetails.Key).To(Equal("test-virtual-maven-repo-read"))
			Expect(repoDetails.Description).To(Equal("Test Virtual Maven Read Repository"))
			Expect(repoDetails.GetRepoType()).To(Equal("virtual"))
			Expect(repoDetails.PackageType).To(Equal("maven"))
		})
	})

	When("a new repository is created with non existing local repo in read instance", func() {
		It("should not exist in Artifactory read instance", func(ctx SpecContext) {
			By("Creating a virtual repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.VirtualMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-virtual-maven-repo-read",
				},
				Spec: v1alpha1.VirtualMavenRepositorySpec{
					ForProvider: v1alpha1.VirtualMavenRepositoryParameters{
						Description:   ptr.To("Test Virtual Maven Read Repository"),
						RepoLayoutRef: ptr.To("maven-2-default"),
						Repositories: []*string{
							ptr.To("test-local-maven-read-repo-nonexistent"),
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
				err := k8sClient.Delete(ctx, &v1alpha1.VirtualMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-virtual-maven-repo-read",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the virtual repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.VirtualMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-virtual-maven-repo-read"}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the virtual repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.VirtualMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-virtual-maven-repo-read"}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionFalse &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionFalse &&
					strings.Contains(repo.Status.GetCondition(v1.TypeSynced).Message,
						"Repository test-local-maven-read-repo-nonexistent does not exist")
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the virtual repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository("test-virtual-maven-repo-read", &repoDetails)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("400"))
			Expect(err.Error()).To(ContainSubstring("Bad Request"))
		})
	})
})
