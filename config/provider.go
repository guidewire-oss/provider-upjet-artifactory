/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	"context"
	_ "embed"

	"github.com/guidewire-oss/provider-jfrogartifactory/config/artifactorygroup"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/artifactoryuser"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/localmavenrepository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/localnpmrepository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/remotemavenrepository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/remotenpmrepository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/repository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/virtualmavenrepository"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/virtualnpmrepository"

	artifactory "github.com/jfrog/terraform-provider-artifactory/v12/pkg/artifactory/provider"

	ujconfig "github.com/crossplane/upjet/pkg/config"
	"github.com/crossplane/upjet/pkg/registry/reference"
)

const (
	resourcePrefix = "jfrogartifactory"
	modulePath     = "github.com/guidewire-oss/provider-jfrogartifactory"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider(ctx context.Context) (*ujconfig.Provider, error) {
	sdkProvider := artifactory.SdkV2()
	fwProvider := artifactory.Framework()

	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("jfrogartifactory.upbound.io"),
		ujconfig.WithShortName("artifactory"),
		ujconfig.WithIncludeList(resourceList(cliReconciledExternalNameConfigs)),
		ujconfig.WithTerraformPluginSDKIncludeList(resourceList(terraformPluginSDKExternalNameConfigs)),
		ujconfig.WithTerraformPluginFrameworkIncludeList(resourceList(terraformPluginFrameworkExternalNameConfigs)),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithReferenceInjectors([]ujconfig.ReferenceInjector{reference.NewInjector(modulePath)}),
		ujconfig.WithTerraformProvider(sdkProvider),
		ujconfig.WithTerraformPluginFrameworkProvider(fwProvider()),
		ujconfig.WithDefaultResourceOptions(
			resourceConfigurator(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		repository.Configure,
		localnpmrepository.Configure,
		remotenpmrepository.Configure,
		virtualnpmrepository.Configure,
		localmavenrepository.Configure,
		remotemavenrepository.Configure,
		virtualmavenrepository.Configure,
		artifactoryuser.Configure,
		artifactorygroup.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc, nil
}

// resourceList returns the list of resources that have external
// name configured in the specified table.
func resourceList(t map[string]ujconfig.ExternalName) []string {
	l := make([]string, len(t))
	i := 0
	for n := range t {
		// Expected format is regex and we'd like to have exact matches.
		l[i] = n + "$"
		i++
	}
	return l
}
