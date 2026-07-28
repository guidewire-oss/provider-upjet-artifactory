/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	"github.com/guidewire-oss/provider-jfrogartifactory/config/artifactorygroup"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/artifactorypermissiontarget"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/artifactoryuser"
	"github.com/guidewire-oss/provider-jfrogartifactory/config/repository"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

const (
	resourcePrefix = "jfrogartifactory"
	modulePath     = "github.com/guidewire-oss/provider-jfrogartifactory"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// resourceConfigurators holds the per-resource custom config functions applied
// by GetProvider. These functions set the API short group and external-name
// behaviour for each resource.
var resourceConfigurators = []func(provider *ujconfig.Provider){
	repository.Configure,
	artifactoryuser.Configure,
	artifactorygroup.Configure,
	artifactorypermissiontarget.Configure,
}

// GetProvider returns the provider configuration. Only cluster-scoped
// resources are generated; the provider does not serve namespaced
// (Crossplane v2) APIs.
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("upbound.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range resourceConfigurators {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
