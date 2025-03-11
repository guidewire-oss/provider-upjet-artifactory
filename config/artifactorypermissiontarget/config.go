package artifactorypermissiontarget

import "github.com/crossplane/upjet/pkg/config"

// Configure the "artifactory_*_repository" resources.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_permission_target", func(r *config.Resource) {
		r.ShortGroup = "permissiontarget"
		r.Kind = "ArtifactoryPermissionTarget"
	})
}
