/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
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
// to both the cluster-scoped and namespaced providers. These functions only set
// scope-independent options (ShortGroup, external-name), so the same list is
// reused by GetProvider and GetProviderNamespaced rather than duplicating the
// config packages under config/cluster and config/namespaced.
var resourceConfigurators = []func(provider *ujconfig.Provider){
	repository.Configure,
	localnpmrepository.Configure,
	remotenpmrepository.Configure,
	virtualnpmrepository.Configure,
	localmavenrepository.Configure,
	remotemavenrepository.Configure,
	virtualmavenrepository.Configure,
	artifactoryuser.Configure,
	artifactorygroup.Configure,
}

// GetProvider returns the cluster-scoped provider configuration.
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

// GetProviderNamespaced returns the namespaced provider configuration
// (Crossplane v2 / Upjet v2). The root group gains the ".m" infix so the
// generated API groups become e.g. jfrogartifactory.m.upbound.io, and generated
// example manifests are namespaced.
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("m.upbound.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range resourceConfigurators {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
