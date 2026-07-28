package remotemavenrepository

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure the "artifactory_*_repository" resources.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_remote_maven_repository", func(r *config.Resource) {
		r.ShortGroup = "repository"
		r.Kind = "RemoteMavenRepository"

		// password_wo is a Terraform write-only attribute (provider 12.11.8+),
		// which requires Terraform 1.11 or later. This provider is pinned to
		// Terraform 1.5.7 because 1.6+ is BSL licensed, so the write-only
		// semantics are unavailable and the attribute cannot behave correctly.
		// Drop it rather than exposing a field that silently misbehaves.
		delete(r.TerraformResource.Schema, "password_wo")
		delete(r.TerraformResource.Schema, "password_wo_version")
	})
}
