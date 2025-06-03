package e2e_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/magefile/mage/sh"

	rt "github.com/jfrog/jfrog-client-go/artifactory"
	rtAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	rtConfig "github.com/jfrog/jfrog-client-go/config"
	v1alpha1group "github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
	v1alpha1repository "github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
	v1alpha1user "github.com/guidewire-oss/provider-jfrogartifactory/apis/jfrogartifactory/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var rtReadClient rt.ArtifactoryServicesManager
var rtWriteClient rt.ArtifactoryServicesManager
var k8sClient client.Client

var _ = ginkgo.SynchronizedBeforeSuite(func() {
	// Applying provider configs
	outsb := strings.Builder{}
	errsb := strings.Builder{}
	outsb.Reset()
	errsb.Reset()
	fmt.Printf("Applying provider configs\n")
	_, err := sh.Exec(nil, &outsb, &errsb, "kubectl", "apply", "-f", "providerconfig-read.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = sh.Exec(nil, &outsb, &errsb, "kubectl", "apply", "-f", "providerconfig-write.yaml")
	Expect(err).NotTo(HaveOccurred())
	fmt.Printf("Applied provider configs\n")
}, func() {
	// Set up the Artifactory client to read instance
	By("Setting up the Artifactory client to read instance")
	ctx, cancel := context.WithCancel(context.Background())
	DeferCleanup(cancel)

	serviceDetails := rtAuth.NewArtifactoryDetails()
	u, err := url.Parse(os.Getenv("READ_URL"))
	Expect(err).NotTo(HaveOccurred())
	u.Path = path.Join(u.Path, "artifactory")
	serviceDetails.SetUrl(u.String())
	serviceDetails.SetUser(os.Getenv("READ_CREDENTIAL_USER"))
	serviceDetails.SetPassword(os.Getenv("READ_CREDENTIAL_ACCESS_TOKEN"))

	serviceConfig, err := rtConfig.NewConfigBuilder().
		SetServiceDetails(serviceDetails).
		SetDryRun(false).
		SetContext(ctx).
		Build()
	Expect(err).NotTo(HaveOccurred())

	rtReadClient, err = rt.New(serviceConfig)
	Expect(err).NotTo(HaveOccurred())

	// Set up the Artifactory client to write instance
	By("Setting up the Artifactory client to write instance")
	ctx, cancel = context.WithCancel(context.Background())
	DeferCleanup(cancel)

	serviceDetails = rtAuth.NewArtifactoryDetails()
	u, err = url.Parse(os.Getenv("WRITE_URL"))
	Expect(err).NotTo(HaveOccurred())
	u.Path = path.Join(u.Path, "artifactory")
	serviceDetails.SetUrl(u.String())
	serviceDetails.SetUser(os.Getenv("WRITE_CREDENTIAL_USER"))
	serviceDetails.SetPassword(os.Getenv("WRITE_CREDENTIAL_ACCESS_TOKEN"))

	serviceConfig, err = rtConfig.NewConfigBuilder().
		SetServiceDetails(serviceDetails).
		SetDryRun(false).
		SetContext(ctx).
		Build()
	Expect(err).NotTo(HaveOccurred())

	rtWriteClient, err = rt.New(serviceConfig)
	Expect(err).NotTo(HaveOccurred())

	// Set up the Kubernetes client
	By("Setting up the Kubernetes client")
	scheme := runtime.NewScheme()
	err = v1alpha1repository.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = v1alpha1user.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = v1alpha1group.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = corev1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	cfg := config.GetConfigOrDie()
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = ginkgo.SynchronizedAfterSuite(func() {}, func() {
	// Deleting provider configs
	outsb := strings.Builder{}
	errsb := strings.Builder{}
	outsb.Reset()
	errsb.Reset()
	fmt.Printf("Deleting provider configs for read and write \n")
	_, err := sh.Exec(nil, &outsb, &errsb, "kubectl", "delete", "-f", "providerconfig-read.yaml")
	Expect(err).NotTo(HaveOccurred())
	_, err = sh.Exec(nil, &outsb, &errsb, "kubectl", "delete", "-f", "providerconfig-write.yaml")
	Expect(err).NotTo(HaveOccurred())
	fmt.Printf("Deleted provider configs for read and write\n")
})
