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

	v1alpha1group "github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
	v1alpha1user "github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
)

var _ = Describe("Artifactory Group", func() {
	var userName string
	var email string

	BeforeEach(func(ctx SpecContext) {
		By("Creating an artifactory user resource in Kubernetes")
		userName = fmt.Sprintf("test-artifactory-user-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
		email = fmt.Sprintf("testartifactoryuser%d_%d@guidewire.com", GinkgoRandomSeed(), GinkgoParallelProcess())
		err := k8sClient.Create(ctx, &v1alpha1user.ArtifactoryUser{
			ObjectMeta: metav1.ObjectMeta{
				Name: userName,
			},
			Spec: v1alpha1user.ArtifactoryUserSpec{
				ForProvider: v1alpha1user.ArtifactoryUserParameters{
					Name:  &userName,
					Email: ptr.To(email),
				},
				ResourceSpec: v1.ResourceSpec{
					ProviderConfigReference: &v1.Reference{
						Name: "my-artifactory-providerconfig-read",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for user to be ready in Kubernetes")
		Eventually(func() bool {
			user := &v1alpha1user.ArtifactoryUser{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: userName}, user)
			Expect(err).NotTo(HaveOccurred())
			return user.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
				user.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
		}, "2m", "5s").Should(BeTrue())

		By("Verifying user exists on Artifactory")
		userDetails := rtServices.UserParams{
			UserDetails: rtServices.User{
				Name:  userName,
				Email: email,
			},
		}
		retrievedUser, err := rtReadClient.GetUser(userDetails)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedUser.Name).To(Equal(userName))
		Expect(retrievedUser.Email).To(Equal(email))
	})

	AfterEach(func(ctx SpecContext) {
		By("Deleting the user resource from Kubernetes")
		err := k8sClient.Delete(ctx, &v1alpha1user.ArtifactoryUser{
			ObjectMeta: metav1.ObjectMeta{
				Name: userName,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the user resource to be deleted")
		Eventually(func() bool {
			user := &v1alpha1user.ArtifactoryUser{}
			err := k8sClient.Get(ctx, client.ObjectKey{Name: userName}, user)
			return errors.IsNotFound(err)
		}, "2m", "5s").Should(BeTrue())

		By("Verifying user does not exist in Artifactory")
		userDetails := rtServices.UserParams{
			UserDetails: rtServices.User{
				Name:  userName,
				Email: email,
			},
		}
		retrievedUser, _ := rtReadClient.GetUser(userDetails)
		Expect(retrievedUser).To(BeNil())
	})

	When("a new group is created", func() {
		It("should exists in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory Group
			groupName := fmt.Sprintf("test-artifactory-group-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			groupDescription := fmt.Sprintf("Test Artifactory Group %d %d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating an artifactory group resource in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1group.ArtifactoryGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: groupName,
				},
				Spec: v1alpha1group.ArtifactoryGroupSpec{
					ForProvider: v1alpha1group.ArtifactoryGroupParameters{
						Description: ptr.To(groupDescription),
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
				By("Deleting the group resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1group.ArtifactoryGroup{
					ObjectMeta: metav1.ObjectMeta{
						Name: groupName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the group resource to be deleted")
				Eventually(func() bool {
					group := &v1alpha1group.ArtifactoryGroup{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: groupName}, group)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())

				By("Verifying group does not exist in Artifactory")
				groupDetails := rtServices.GroupParams{
					GroupDetails: rtServices.Group{
						Name:        groupName,
						Description: groupDescription,
					},
				}
				retrievedGroup, _ := rtReadClient.GetGroup(groupDetails)
				Expect(retrievedGroup).To(BeNil())
			})

			// Check if object is ready
			By("Waiting for group to be ready in Kubernetes")
			Eventually(func() bool {
				group := &v1alpha1group.ArtifactoryGroup{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: groupName}, group)
				Expect(err).NotTo(HaveOccurred())
				return group.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					group.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			// Verify if group is created on Artifactory
			By("Verifying group exists on Artifactory")
			groupDetails := rtServices.GroupParams{
				GroupDetails: rtServices.Group{
					Name:        groupName,
					Description: groupDescription,
				},
			}
			retrievedGroup, err := rtReadClient.GetGroup(groupDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedGroup.Name).To(Equal(groupName))
			Expect(retrievedGroup.Description).To(Equal(groupDescription))
		})
	})

	When("a new group is created with user added", func() {
		It("should exists in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory Group
			groupName := fmt.Sprintf("test-artifactory-group-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			groupDescription := fmt.Sprintf("Test Artifactory Group %d %d", GinkgoRandomSeed(), GinkgoParallelProcess())
			By("Creating an artifactory group with user resource in Kubernetes")
			err := k8sClient.Create(ctx, &v1alpha1group.ArtifactoryGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: groupName,
				},
				Spec: v1alpha1group.ArtifactoryGroupSpec{
					ForProvider: v1alpha1group.ArtifactoryGroupParameters{
						Description: ptr.To(groupDescription),
						UsersNames: []*string{
							ptr.To(userName),
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
				By("Deleting the group resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1group.ArtifactoryGroup{
					ObjectMeta: metav1.ObjectMeta{
						Name: groupName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the group resource to be deleted")
				Eventually(func() bool {
					group := &v1alpha1group.ArtifactoryGroup{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: groupName}, group)
					return errors.IsNotFound(err)
				}, "2m", "5s").Should(BeTrue())

				By("Verifying group does not exist in Artifactory")
				groupDetails := rtServices.GroupParams{
					GroupDetails: rtServices.Group{
						Name:        groupName,
						Description: groupDescription,
					},
				}
				retrievedGroup, _ := rtReadClient.GetGroup(groupDetails)
				Expect(retrievedGroup).To(BeNil())
			})

			// Check if object is ready
			By("Waiting for group to be ready in Kubernetes")
			Eventually(func() bool {
				group := &v1alpha1group.ArtifactoryGroup{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: groupName}, group)
				Expect(err).NotTo(HaveOccurred())
				return group.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					group.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			// Verify if group is created on Artifactory
			By("Verifying group exists on Artifactory")
			groupDetails := rtServices.GroupParams{
				GroupDetails: rtServices.Group{
					Name:        groupName,
					Description: groupDescription,
				},
				IncludeUsers: true,
			}
			retrievedGroup, err := rtReadClient.GetGroup(groupDetails)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedGroup.Name).To(Equal(groupName))
			Expect(retrievedGroup.Description).To(Equal(groupDescription))
			Expect(retrievedGroup.UsersNames).To(ContainElement(userName))
		})
	})
})
