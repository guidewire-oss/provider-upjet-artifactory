package artifactorygroup

import "github.com/crossplane/upjet/pkg/config"

// Configure the "artifactory_*" resources.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_group", func(r *config.Resource) {
		r.Kind = "ArtifactoryGroup"
		r.ShortGroup = "jfrogartifactory"
	})
}
