//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
	"github.com/myorg/provider-jfrogartifactory/e2e"
)

// SetupE2E sets up the environment for end-to-end tests.
func SetupE2E() error {
	err := e2e.EnsureKindCluster("kind")

	if err != nil {
		return err
	}

	// Uncomment if you want to install artifactory in a pod in a kind cluster
	//err = e2e.EnsureArtifactory()
	//if err != nil {
	//	return err
	//}

	return e2e.UpdateCredentials()
}

// TestE2E runs the end-to-end tests.
func TestE2E() error {
	// See: https://onsi.github.io/ginkgo/#recommended-continuous-integration-configuration
	return sh.RunV("ginkgo", "-r", "-v", "-p",
		"--fail-on-pending",
		"--randomize-all",
		"--randomize-suites",
		"--keep-going",
		"--procs=2",
		"e2e",
	)
}

// Lint runs golangci-lint on the project.
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}
