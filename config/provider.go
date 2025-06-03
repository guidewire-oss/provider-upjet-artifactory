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

	ujconfig "github.com/crossplane/upjet/pkg/config"
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
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("upbound.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
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
	return pc
}
