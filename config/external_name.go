/*
Copyright 2022 Upbound Inc.
*/

package config

import (
	"errors"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// keyAsIdentifier returns an external name configuration for repository
// resources whose Terraform schema no longer exposes an "id" attribute.
//
// The JFrog provider migrated its local and remote repository resources to the
// Terraform Plugin Framework (12.7.0, 12.8.2, 12.8.3), which dropped the
// synthetic "id" that the SDKv2 implementations carried. The default
// GetExternalNameFn (config.IDAsExternalName) reads tfstate["id"] and would
// fail every Observe with "cannot find id in tfstate". These resources are
// identified by their repository key instead, which is always present in the
// state because it is the required identifier argument.
func keyAsIdentifier() config.ExternalName {
	e := config.ParameterAsIdentifier("key")
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		if key, ok := tfstate["key"].(string); ok && key != "" {
			return key, nil
		}
		return "", errors.New("cannot find key in tfstate")
	}
	return e
}

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Import requires using a randomly generated ID from provider: nl-2e21sda
	// TODO: Not implemented yet: "artifactory_unmanaged_user":           config.NameAsIdentifier,

	// Plugin Framework resources: no "id" in the schema, keyed on "key".
	"artifactory_local_generic_repository": keyAsIdentifier(),
	"artifactory_local_npm_repository":     keyAsIdentifier(),
	"artifactory_local_maven_repository":   keyAsIdentifier(),
	"artifactory_remote_npm_repository":    keyAsIdentifier(),
	"artifactory_remote_maven_repository":  keyAsIdentifier(),

	// Still SDKv2 upstream: "id" is present in the schema.
	"artifactory_virtual_npm_repository":   config.ParameterAsIdentifier("key"),
	"artifactory_virtual_maven_repository": config.ParameterAsIdentifier("key"),
	"artifactory_user":                     config.ParameterAsIdentifier("name"),
	"artifactory_group":                    config.ParameterAsIdentifier("name"),
	"artifactory_permission_target":        config.ParameterAsIdentifier("name"),
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
