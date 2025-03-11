package e2e

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
	"maps"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"

	v1alpha1permission "github.com/myorg/provider-jfrogartifactory/apis/permissiontarget/v1alpha1"
)

var _ = Describe("Artifactory Permission Target", func() {
	var userName string
	var email string
	var groupName string
	var groupDescription string
	var repoName string

	BeforeEach(func(ctx SpecContext) {
		// ------------- Create Artifactory User ---------------------
		By("Creating an artifactory user resource in Artifactory")
		userName = fmt.Sprintf("test-artifactory-user-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		email = fmt.Sprintf("testartifactoryuser%d_%d@guidewire.com", GinkgoRandomSeed(), GinkgoParallelProcess())
		userParams := rtServices.UserParams{
			UserDetails: rtServices.User{
				Name:     userName,
				Email:    email,
				Password: "Testpassword1",
			},
		}
		err := rtReadClient.CreateUser(userParams)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx SpecContext) {
			err := rtReadClient.DeleteUser(userName)
			Expect(err).NotTo(HaveOccurred())
		})
		// --------------------------------------------------------------

		// ------------- Create Artifactory Group ---------------------
		By("Creating an artifactory group resource in Artifactory")
		groupName = fmt.Sprintf("test-artifactory-group-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		groupDescription = fmt.Sprintf("Test Artifactory Group %d %d", GinkgoRandomSeed(), GinkgoParallelProcess())
		groupParams := rtServices.GroupParams{
			GroupDetails: rtServices.Group{
				Name:        groupName,
				Description: groupDescription,
			},
		}
		err = rtReadClient.CreateGroup(groupParams)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx SpecContext) {
			err := rtReadClient.DeleteGroup(groupName)
			Expect(err).NotTo(HaveOccurred())
		})
		// --------------------------------------------------------------

		// ------------- Create Local Maven Repo ---------------------
		By("Creating an artifactory local maven repository resource in Artifactory")
		repoName = fmt.Sprintf("test-local-maven-repo-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		params := rtServices.NewMavenLocalRepositoryParams()
		params.Key = repoName
		err = rtReadClient.CreateLocalRepository().Maven(params)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx SpecContext) {
			err := rtReadClient.DeleteRepository(repoName)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("a new permission target is created", func() {
		It("should exists in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory Permission Target
			permissionName := fmt.Sprintf("test-artifactory-permission-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating an artifactory permission target resource in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1permission.ArtifactoryPermissionTarget{
				ObjectMeta: metav1.ObjectMeta{
					Name: permissionName,
				},
				Spec: v1alpha1permission.ArtifactoryPermissionTargetSpec{
					ForProvider: v1alpha1permission.ArtifactoryPermissionTargetParameters{
						Repo: []v1alpha1permission.RepoParameters{
							{
								Repositories: []*string{
									&repoName,
								},
								Actions: []v1alpha1permission.RepoActionsParameters{
									{
										Users: []v1alpha1permission.RepoActionsUsersParameters{
											{
												Name: &userName,
												Permissions: []*string{
													ptr.To("read"),
												},
											},
										},
										Groups: []v1alpha1permission.RepoActionsGroupsParameters{
											{
												Name: &groupName,
												Permissions: []*string{
													ptr.To("read"),
												},
											},
										},
									},
								},
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
				By("Deleting the permission target resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1permission.ArtifactoryPermissionTarget{
					ObjectMeta: metav1.ObjectMeta{
						Name: permissionName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the permission target resource to be deleted")
				Eventually(func() bool {
					permission := &v1alpha1permission.ArtifactoryPermissionTarget{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: permissionName}, permission)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())

				By("Verifying permission target does not exist in Artifactory")
				retrievedPermissionTarget, _ := rtReadClient.GetPermissionTarget(permissionName)
				Expect(retrievedPermissionTarget).To(BeNil())
			})

			// Check if object is ready
			By("Waiting for permission target to be ready in Kubernetes")
			Eventually(func() bool {
				permission := &v1alpha1permission.ArtifactoryPermissionTarget{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: permissionName}, permission)
				Expect(err).NotTo(HaveOccurred())
				return permission.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					permission.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			// Verify if permission target is created on Artifactory
			By("Verifying permission target exists on Artifactory")
			retrievedPermissionTarget, err := rtReadClient.GetPermissionTarget(permissionName)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedPermissionTarget.Name).To(Equal(permissionName))
			Expect(retrievedPermissionTarget.Repo.Repositories).To(ContainElement(repoName))
			Expect(slices.Collect(maps.Keys(retrievedPermissionTarget.Repo.Actions.Users))).To(ContainElement(userName))
			Expect(retrievedPermissionTarget.Repo.Actions.Users[userName]).To(ContainElement("read"))
			Expect(retrievedPermissionTarget.Repo.Actions.Groups[groupName]).To(ContainElement("read"))
		})
	})
})
