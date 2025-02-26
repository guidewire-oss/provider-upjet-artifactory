package e2e_test

import (
	"fmt"
	"os"
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

var _ = Describe("RemoteMavenRepository", func() {
	var localRepoName string

	BeforeEach(func(ctx SpecContext) {
		localRepoName = fmt.Sprintf("test-local-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		By("Creating a local repository resource with write ProviderConfig in Kubernetes")
		err := k8sClient.Create(ctx, &v1alpha1.LocalMavenRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: localRepoName,
			},
			Spec: v1alpha1.LocalMavenRepositorySpec{
				ForProvider: v1alpha1.LocalMavenRepositoryParameters{
					Description: ptr.To("Test Local Maven Write Repository"),
				},
				ResourceSpec: v1.ResourceSpec{
					ProviderConfigReference: &v1.Reference{
						Name: "my-artifactory-providerconfig-write",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the local repository to be ready in Kubernetes")
		Eventually(func() bool {
			repo := &v1alpha1.LocalMavenRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: localRepoName}, repo)
			Expect(err).NotTo(HaveOccurred())
			return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
				repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
		}, "2m", "5s").Should(BeTrue())

		// Test for actual resource existence in artifactory instance
		By("Verifying the repository exists in Artifactory write instances")
		repoDetails := rtServices.RepositoryDetails{}
		err = rtWriteClient.GetRepository(localRepoName, &repoDetails)
		Expect(err).NotTo(HaveOccurred())
		Expect(repoDetails.Key).To(Equal(localRepoName))
		Expect(repoDetails.Description).To(Equal("Test Local Maven Write Repository"))
		Expect(repoDetails.GetRepoType()).To(Equal("local"))
		Expect(repoDetails.PackageType).To(Equal("maven"))

	})

	AfterEach(func(ctx SpecContext) {
		By("Deleting the local repository resource from Kubernetes")
		err := k8sClient.Delete(ctx, &v1alpha1.LocalMavenRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: localRepoName,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the local repository resource to be deleted")
		Eventually(func() bool {
			repo := &v1alpha1.LocalMavenRepository{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: localRepoName}, repo)
			return errors.IsNotFound(err)
		}, "2m", "5s").Should(BeTrue())
		// 	Test actual repository to be deleted
		By("Verifying the repository does not exists in Artifactory write instance")
		repoDetails := rtServices.RepositoryDetails{}
		err = rtWriteClient.GetRepository(localRepoName, &repoDetails)
		Expect(err).To(HaveOccurred())
	})

	When("a new maven repository is created with valid remote artifactory instance creds and pointing to a valid local repo in remote instance", func() {
		It("should exist in Artifactory read instance", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-remote-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a remote repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.RemoteMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.RemoteMavenRepositorySpec{
					ForProvider: v1alpha1.RemoteMavenRepositoryParameters{
						Description: ptr.To("Test Remote Maven Repository Read"),
						// Replace the below from git history
						URL: ptr.To(os.Getenv("WRITE_URL") + "/artifactory/" + localRepoName),
						ContentSynchronisation: []v1alpha1.ContentSynchronisationParameters{
							{
								Enabled:                      ptr.To(true),
								PropertiesEnabled:            ptr.To(true),
								SourceOriginAbsenceDetection: ptr.To(true),
								StatisticsEnabled:            ptr.To(true),
							},
						},
						Username: ptr.To(os.Getenv("WRITE_CREDENTIAL_USER")),
						PasswordSecretRef: &v1.SecretKeySelector{
							Key: "passwords",
							SecretReference: v1.SecretReference{
								Name:      "secretremote",
								Namespace: "default",
							},
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
				By("Deleting the remote repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.RemoteMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the remote repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.RemoteMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the remote repository to be ready in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.RemoteMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the repository exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository(repoName, &repoDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(repoDetails.Key).To(Equal(repoName))
			Expect(repoDetails.Description).To(Equal("Test Remote Maven Repository Read"))
			Expect(repoDetails.GetRepoType()).To(Equal("remote"))
			Expect(repoDetails.PackageType).To(Equal("maven"))
		})
	})

	When("a new repository is created invalid creds for remote artifactory instance", func() {
		It("should not exist in Artifactory read instance", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-remote-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a remote repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.RemoteMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.RemoteMavenRepositorySpec{
					ForProvider: v1alpha1.RemoteMavenRepositoryParameters{
						Description: ptr.To("Test Remote Maven Repository Read"),
						// Replace the below from git history
						URL: ptr.To(os.Getenv("WRITE_URL") + "/artifactory/" + localRepoName),
						ContentSynchronisation: []v1alpha1.ContentSynchronisationParameters{
							{
								Enabled:                      ptr.To(true),
								PropertiesEnabled:            ptr.To(true),
								SourceOriginAbsenceDetection: ptr.To(true),
								StatisticsEnabled:            ptr.To(true),
							},
						},
						Username: ptr.To(os.Getenv("WRITE_CREDENTIAL_USER")),
						PasswordSecretRef: &v1.SecretKeySelector{
							Key: "invalid-passwords",
							SecretReference: v1.SecretReference{
								Name:      "secretremote",
								Namespace: "default",
							},
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
				By("Deleting the remote repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.RemoteMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the remote repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.RemoteMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the remote repository to fail in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.RemoteMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionFalse &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionFalse &&
					strings.Contains(repo.Status.GetCondition(v1.TypeSynced).Message,
						"Upstream repository credentials or user permissions insufficient")
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the repository did not exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository(repoName, &repoDetails)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("400"))
			Expect(err.Error()).To(ContainSubstring("Bad Request"))
		})
	})

	When("a new repository is created pointing to a invalid local repo in remote instance", func() {
		It("should not exist in Artifactory read instance", func(ctx SpecContext) {
			repoName := fmt.Sprintf("test-remote-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating a remote repository resource with read ProviderConfig in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1.RemoteMavenRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name: repoName,
				},
				Spec: v1alpha1.RemoteMavenRepositorySpec{
					ForProvider: v1alpha1.RemoteMavenRepositoryParameters{
						Description: ptr.To("Test Remote Maven Repository Read"),
						URL:         ptr.To(os.Getenv("WRITE_URL") + "/artifactory/test-maven-write-repo/"),
						ContentSynchronisation: []v1alpha1.ContentSynchronisationParameters{
							{
								Enabled:                      ptr.To(true),
								PropertiesEnabled:            ptr.To(true),
								SourceOriginAbsenceDetection: ptr.To(true),
								StatisticsEnabled:            ptr.To(true),
							},
						},
						Username: ptr.To(os.Getenv("WRITE_CREDENTIAL_USER")),
						PasswordSecretRef: &v1.SecretKeySelector{
							Key: "passwords",
							SecretReference: v1.SecretReference{
								Name:      "secretremote",
								Namespace: "default",
							},
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
				By("Deleting the remote repository resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.RemoteMavenRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name: repoName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the remote repository resource to be deleted")
				Eventually(func() bool {
					repo := &v1alpha1.RemoteMavenRepository{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())
			})

			By("Waiting for the remote repository to fail in Kubernetes")
			Eventually(func() bool {
				repo := &v1alpha1.RemoteMavenRepository{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: repoName}, repo)
				Expect(err).NotTo(HaveOccurred())
				return repo.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionFalse &&
					repo.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionFalse &&
					strings.Contains(repo.Status.GetCondition(v1.TypeSynced).Message,
						"Upstream repository not found or not accessible")
			}, "2m", "5s").Should(BeTrue())

			By("Verifying the repository did not exists in Artifactory")
			repoDetails := rtServices.RepositoryDetails{}
			err = rtReadClient.GetRepository(repoName, &repoDetails)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("400"))
			Expect(err.Error()).To(ContainSubstring("Bad Request"))
		})
	})

})
