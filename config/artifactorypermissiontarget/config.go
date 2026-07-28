package artifactorypermissiontarget

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure the "artifactory_permission_target" resource.
//
// Note: the upstream Terraform resource has been deprecated since provider
// 10.3.1 in favour of the platform_permission resource in the JFrog Platform
// provider. It remains fully functional and is still present in the schema, but
// it will eventually need to migrate to a separate provider package, since
// upjet binds one Terraform provider per Crossplane provider.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("artifactory_permission_target", func(r *config.Resource) {
		r.ShortGroup = "permission"
		r.Kind = "ArtifactoryPermissionTarget"
	})
}
