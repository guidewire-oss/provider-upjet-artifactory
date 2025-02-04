package remotemavenrepository

import "github.com/crossplane/upjet/pkg/config"

// Configure the "artifactory_*_repository" resources.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_remote_maven_repository", func(r *config.Resource) {
		r.ShortGroup = "repository"
		r.Kind = "RemoteMavenRepository"
	})
}
