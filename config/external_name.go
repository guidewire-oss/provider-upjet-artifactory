/*
Copyright 2022 Upbound Inc.
*/

package config

import "github.com/crossplane/upjet/pkg/config"

// TerraformPluginFrameworkExternalNameConfigs contains all external
// name configurations belonging to Terraform Plugin Framework
// resources to be reconciled under the no-fork architecture for this
// provider.
var terraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{

	"artifactory_group":                    config.ParameterAsIdentifier("name"),
	"artifactory_local_generic_repository": config.ParameterAsIdentifier("key"),
	"artifactory_local_maven_repository":   config.ParameterAsIdentifier("key"),
	"artifactory_local_npm_repository":     config.ParameterAsIdentifier("key"),
	"artifactory_remote_maven_repository":  config.ParameterAsIdentifier("key"),
	"artifactory_remote_npm_repository":    config.ParameterAsIdentifier("key"),
	"artifactory_user":                     config.ParameterAsIdentifier("name"),
}

// TerraformPluginSDKExternalNameConfigs contains all external name configurations
// belonging to Terraform Plugin SDKv2 resources to be reconciled
// under the no-fork architecture for this provider.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{
	// Import requires using a randomly generated ID from provider: nl-2e21sda
	// TODO: Not implemented yet: "artifactory_unmanaged_user":           config.NameAsIdentifier,
	"artifactory_virtual_maven_repository": config.ParameterAsIdentifier("key"),
	"artifactory_virtual_npm_repository":   config.ParameterAsIdentifier("key"),
}

// cliReconciledExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the CLI-based
// architecture for this provider.
var cliReconciledExternalNameConfigs = map[string]config.ExternalName{}

// resourceConfigurator applies all external name configs
// listed in the table terraformPluginSDKExternalNameConfigs and
// cliReconciledExternalNameConfigs and sets the version
// of those resources to v1beta1. For those resource in
// terraformPluginSDKExternalNameConfigs, it also sets
// config.Resource.UseNoForkClient to `true`.
func resourceConfigurator() config.ResourceOption {
	return func(r *config.Resource) {
		// if configured both for the no-fork and CLI based architectures,
		// no-fork configuration prevails
		e, configured := terraformPluginSDKExternalNameConfigs[r.Name]
		if !configured {
			e, configured = cliReconciledExternalNameConfigs[r.Name]
		}
		if !configured {
			return
		}
		r.Version = "v1alpha1"
		r.ExternalName = e
	}
}
