package artifactoryuser

import "github.com/crossplane/upjet/pkg/config"

// Configure the "artifactory_*" resources.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_user", func(r *config.Resource) {
		r.Kind = "ArtifactoryUser"
		r.ShortGroup = "user"
		r.ExternalName.OmittedFields = []string{
			"password_policy",
		}
	})
}
