package e2e_test

import (
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	rtServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/myorg/provider-jfrogartifactory/apis/user/v1alpha1"
)

var _ = Describe("Artifactory User", func() {

	When("a new user is created", func() {
		It("should exists in Artifactory read instance", func(ctx SpecContext) {
			// Create Kubernetes object of Artifactory User
			userName := fmt.Sprintf("test-artifactory-user-%d-%d", GinkgoRandomSeed(), GinkgoParallelProcess())
			email := fmt.Sprintf("testartifactoryuser%d_%d@guidewire.com", GinkgoRandomSeed(), GinkgoParallelProcess())
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "password",
					Namespace: "default",
				},
				Type: "Opaque",
				Data: map[string][]byte{
					"password": ([]byte(b64.StdEncoding.EncodeToString([]byte(generateRandomPasswordString(10))))),
				},
			}
			By("Creating an artifactory user resource in Kubernetes")
			err := k8sClient.Create(ctx, secret)
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func(ctx SpecContext) {
				By("Deleting the secret resource from Kubernetes")
				err := k8sClient.Delete(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "password",
						Namespace: "default",
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})
			err = k8sClient.Create(ctx, &v1alpha1.ArtifactoryUser{
				ObjectMeta: metav1.ObjectMeta{
					Name: userName,
				},
				Spec: v1alpha1.ArtifactoryUserSpec{
					ForProvider: v1alpha1.ArtifactoryUserParameters{
						Name:  &userName,
						Email: ptr.To(email),
						PasswordSecretRef: &v1.SecretKeySelector{
							SecretReference: v1.SecretReference{
								Name:      "password",
								Namespace: "default",
							},
							Key: "password",
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
				By("Deleting the user resource from Kubernetes")
				err := k8sClient.Delete(ctx, &v1alpha1.ArtifactoryUser{
					ObjectMeta: metav1.ObjectMeta{
						Name: userName,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for the user resource to be deleted")
				Eventually(func() bool {
					user := &v1alpha1.ArtifactoryUser{}
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

			// Check if object is ready
			By("Waiting for user to be ready in Kubernetes")
			Eventually(func() bool {
				user := &v1alpha1.ArtifactoryUser{}
				err := k8sClient.Get(ctx, client.ObjectKey{Name: userName}, user)
				Expect(err).NotTo(HaveOccurred())
				return user.Status.GetCondition(v1.TypeReady).Status == corev1.ConditionTrue &&
					user.Status.GetCondition(v1.TypeSynced).Status == corev1.ConditionTrue
			}, "2m", "5s").Should(BeTrue())

			// Verify if user is created on Artifactory
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
	})
})

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return strconv.Itoa(rand.Int()) + string(b)
}

func generateRandomPasswordString(length int) string {
	return StringWithCharset(length, charset)
}
