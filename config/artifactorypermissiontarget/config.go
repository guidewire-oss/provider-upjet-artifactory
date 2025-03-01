package artifactorypermissiontarget

import "github.com/crossplane/upjet/pkg/config"

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_permission_target", func(r *config.Resource) {
		r.Kind = "ArtifactoryPermissionTarget"
		r.ShortGroup = "permissiontarget"
	})
}