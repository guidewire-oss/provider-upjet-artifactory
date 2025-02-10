package e2e_test

import (
	"context"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	rt "github.com/jfrog/jfrog-client-go/artifactory"
	rtAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	rtServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	rtConfig "github.com/jfrog/jfrog-client-go/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/myorg/provider-jfrogartifactory/apis/repository/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("E2E Tests", func() {
	var rtClient rt.ArtifactoryServicesManager
	var k8sClient client.Client

	BeforeEach(func() {
		By("Setting up the Artifactory client")
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)

		serviceDetails := rtAuth.NewArtifactoryDetails()
		serviceDetails.SetUrl("url")
		serviceDetails.SetUser("user")
		serviceDetails.SetPassword("pass")

		serviceConfig, err := rtConfig.NewConfigBuilder().
			SetServiceDetails(serviceDetails).
			SetDryRun(false).
			SetContext(ctx).
			Build()
		Expect(err).NotTo(HaveOccurred())

		rtClient, err = rt.New(serviceConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	// Set up the Kubernetes client
	BeforeEach(func() {
		By("Setting up the Kubernetes client")
		scheme := runtime.NewScheme()
		err := v1alpha1.AddToScheme(scheme)
		Expect(err).NotTo(HaveOccurred())

		cfg := config.GetConfigOrDie()
		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).NotTo(BeNil())
	})

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
				err = rtClient.GetRepository("test-repo", &repoDetails)
				Expect(err).NotTo(HaveOccurred())
				Expect(repoDetails.Key).To(Equal("test-repo"))
				Expect(repoDetails.Description).To(Equal("Test repository"))
				Expect(repoDetails.GetRepoType()).To(Equal("local"))
				Expect(repoDetails.PackageType).To(Equal("generic"))
			})
		})
	})
})
