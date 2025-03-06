package artifactorypermissiontarget

import "github.com/crossplane/upjet/pkg/config"

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("platform_permission", func(r *config.Resource) {
		r.Kind = "ArtifactoryPermissionTarget"
		r.ShortGroup = "permissiontarget"
	})
}