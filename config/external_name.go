/*
Copyright 2022 Upbound Inc.
*/

package config

import (
	"errors"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/guidewire-oss/provider-jfrogartifactory/config/repository"
)

// ExternalNameConfigs contains the external name configurations for the
// non-repository resources. Repositories are handled by
// repositoryExternalName, which derives the configuration from the schema.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Import requires using a randomly generated ID from provider: nl-2e21sda
	// TODO: Not implemented yet: "artifactory_unmanaged_user":           config.NameAsIdentifier,
	"artifactory_user":              config.ParameterAsIdentifier("name"),
	"artifactory_group":             config.ParameterAsIdentifier("name"),
	"artifactory_permission_target": config.ParameterAsIdentifier("name"),
}

// repositoryExternalName returns the external name configuration for a
// repository resource. Every repository is identified by its "key" argument.
//
// Resources that JFrog has migrated to the Terraform Plugin Framework no longer
// expose a synthetic "id" attribute, so the default GetExternalNameFn
// (config.IDAsExternalName) reads tfstate["id"] and fails every Observe with
// "cannot find id in tfstate". Rather than maintaining a hand-written list of
// which resources have migrated -- 80 of 111 as of provider 12.11.9, and the
// split shifts with every release -- the decision is taken from the schema.
func repositoryExternalName(r *config.Resource) config.ExternalName {
	e := config.ParameterAsIdentifier("key")
	if r.TerraformResource == nil {
		return e
	}
	if _, hasID := r.TerraformResource.Schema["id"]; hasID {
		return e
	}
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		if key, ok := tfstate["key"].(string); ok && key != "" {
			return key, nil
		}
		return "", errors.New("cannot find key in tfstate")
	}
	return e
}

// ExternalNameConfigurations applies the external name configuration for every
// resource this provider exposes.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
			return
		}
		if repository.Matches(r.Name) {
			r.ExternalName = repositoryExternalName(r)
		}
	}
}

// ExternalNameConfigured returns the include list of resources this provider
// generates: every supported repository, plus the explicitly configured
// resources above.
func ExternalNameConfigured() []string {
	l := make([]string, 0, len(ExternalNameConfigs)+1)
	l = append(l, repository.NameRegex)
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l = append(l, name+"$")
	}
	return l
}
