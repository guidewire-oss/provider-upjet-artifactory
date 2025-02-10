package e2e

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
)

const (
	artifactoryNamespace   = "jfrog"
	artifactoryHelmRelease = "artifactory"
)

// KindClusterExists checks if a Kind cluster with the specified name exists.
func KindClusterExists(name string) (bool, error) {
	output, err := sh.Output("kind", "get", "clusters")

	if err != nil {
		return false, err
	}

	// split the output into an array of lines
	clusters := strings.Split(output, "\n")

	// check if the cluster name is in the list of clusters
	for _, cluster := range clusters {
		if cluster == name {
			return true, nil
		}
	}

	return false, nil
}

// CreateKindCluster creates a new Kind cluster with the given name.
func CreateKindCluster(name string) error {
	err := sh.RunV("kind", "create", "cluster", "--name", name)

	if err != nil {
		return err
	}

	return nil
}

// EnsureKindCluster checks if a Kind cluster with the specified name exists,
// creating it if it does not. It returns an error if the validation or creation fails.
func EnsureKindCluster(name string) error {
	exists, err := KindClusterExists(name)

	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Kind cluster %s already exists.\n", name)
		return nil
	}

	return CreateKindCluster(name)
}

// InstallArtifactory installs the Artifactory Helm release, waits for the pod to be
// created and reach the Ready state within specified timeouts, and returns an error
// if any step fails or times out.
func InstallArtifactory() error {
	const podCreationTimeout = 30 * time.Second
	const podReadyTimeout = 5 * time.Minute
	start := time.Now()

	fmt.Println("Installing Artifactory...")
	err := sh.Run(
		"helm",
		"install",
		artifactoryHelmRelease,
		"--create-namespace",
		"-n",
		artifactoryNamespace,
		"--set",
		"xray.enabled=false",
		"--set",
		os.ExpandEnv("artifactory.artifactory.license.licenseKey=${ARTIFACTORY_LICENSE_KEY}"),
		"--set",
		"artifactory.access.accessConfig.token.allow-basic-auth=true",
		"jfrog/jfrog-platform",
	)

	if err != nil {
		return err
	}

	fmt.Println("Waiting for Artifactory to be ready...")
	for {
		outsb := strings.Builder{}
		errsb := strings.Builder{}

		_, err := sh.Exec(
			nil,
			&outsb,
			&errsb,
			"kubectl",
			"wait",
			"--for=condition=Ready",
			"-n",
			artifactoryNamespace,
			"pod/"+artifactoryHelmRelease+"-0",
		)

		if err != nil {
			if strings.Contains(errsb.String(), "pods \""+artifactoryHelmRelease+"-0\" not found") {
				if time.Since(start) < podCreationTimeout {
					continue
				}

				return fmt.Errorf("artifactory pod not found after %v; will not retry", podCreationTimeout)
			}

			if strings.Contains(errsb.String(), "timed out waiting for the condition on pods/"+artifactoryHelmRelease+"-0") {
				if time.Since(start) < podReadyTimeout {
					continue
				}

				return fmt.Errorf("artifactory pod not ready after %v; will not retry", podReadyTimeout)
			}

			fmt.Printf("Unhandled error: %s\n", err.Error())
			fmt.Printf("Standard output: %s\n", outsb.String())
			fmt.Printf("Error output: %s\n", errsb.String())

			return err
		}

		fmt.Println("Artifactory is ready.")
		return nil
	}
}

// ArtifactoryExists checks if the configured Artifactory Helm release is installed
// in the given namespace. It returns a boolean indicating whether the release is
// found, and an error if the Helm command fails.
func ArtifactoryExists() (bool, error) {
	output, err := sh.Output("helm", "list", "-n", artifactoryNamespace)

	if err != nil {
		return false, err
	}

	// split the output into an array of lines
	releaseInfos := strings.Split(output, "\n")

	// check if the release name is in the list of releases
	for _, releaseInfo := range releaseInfos {
		if strings.Contains(releaseInfo, artifactoryHelmRelease) {
			return true, nil
		}
	}

	return false, nil
}

// EnsureArtifactory verifies the presence of an Artifactory instance and installs
// it if necessary. If Artifactory is already installed, this function does nothing.
func EnsureArtifactory() error {
	exists, err := ArtifactoryExists()

	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Artifactory is already installed.")
		return nil
	}

	return InstallArtifactory()
}

// UpdateCredentials re-runs the credentials job in the Kubernetes cluster to
// generate the Artifactory credentials.
//
// Returns an error if either the deletion (when job exists) or creation fails,
// nil otherwise.
func UpdateCredentials() error {
	outsb := strings.Builder{}
	errsb := strings.Builder{}
	// artifactory_url := os.Getenv("ARTIFACTORY_URL")
	// artifactory_access_token := os.Getenv("ARTIFACTORY_ACCESS_TOKEN")
	//_, err := sh.Exec(nil, &outsb, &errsb, "kubectl", "delete", "job", "create-credentials")

	// if err != nil {
	// 	if !strings.Contains(errsb.String(), "jobs.batch \"create-credentials\" not found") {
	// 		fmt.Printf("Unhandled error: %s\n", err.Error())
	// 		fmt.Printf("Standard output: %s\n", outsb.String())
	// 		fmt.Printf("Error output: %s\n", errsb.String())

	// 		return err
	// 	}
	// }

	outsb.Reset()
	errsb.Reset()

	//_, err = sh.Exec(nil, &outsb, &errsb, "kubectl", "create", "secret", "generic", "artifactory-creds", fmt.Sprintf("--from-literal=ARTIFACTORY_URL=%s", artifactory_url), fmt.Sprintf("--from-literal=ARTIFACTORY_ACCESS_TOKEN=%s", artifactory_access_token))
	//_, err = sh.Exec(nil, &outsb, &errsb, "kubectl", "apply", "-f", "e2e/create-credentials.yaml")

	_, err := sh.Exec(nil, &outsb, &errsb, "kubectl", "create", "secret", "generic", "artifactory-credentials", "--from-file=e2e/creds.json")

	if err != nil {
		fmt.Printf("Unhandled error: %s\n", err.Error())
		fmt.Printf("Standard output: %s\n", outsb.String())
		fmt.Printf("Error output: %s\n", errsb.String())

		return err
	}

	return nil
}
