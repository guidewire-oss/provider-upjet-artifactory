package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/name"
)

// SupportedPackageTypes lists the Artifactory package types this provider
// exposes, as they appear in the Terraform resource names. The Artifactory
// provider offers many more; this list is the deliberate supported subset.
//
// Some package types do not exist for every repository class, which is why the
// tokens are matched against all three classes rather than enumerated per
// class. Non-existent combinations simply match nothing:
//
//   - docker is remote/virtual only; local splits into docker_v1 and docker_v2,
//     and only v2 is supported (v1 is the deprecated legacy registry format).
//   - terraform is remote/virtual only; local splits into terraform_module,
//     terraform_provider and terraformbackend.
//   - helmoci is the OCI-based Helm variant and exists for all three classes.
var SupportedPackageTypes = []string{
	"docker",
	"docker_v2",
	"generic",
	"gradle",
	"helm",
	"helmoci",
	"maven",
	"npm",
	"pypi",
	"terraform",
	"terraform_module",
	"terraform_provider",
	"terraformbackend",
}

// NameRegex matches every supported local, remote and virtual repository
// resource. It is also used as the provider's include list, so that the set of
// generated repositories and the set of configured repositories cannot drift
// apart. Federated repositories are deliberately excluded.
var NameRegex = fmt.Sprintf(
	`^artifactory_(local|remote|virtual)_(%s)_repository$`,
	strings.Join(SupportedPackageTypes, "|"),
)

var nameRe = regexp.MustCompile(NameRegex)

// packageTypeKind maps a Terraform package-type token to the form used in the
// CRD Kind. Upjet's snake-to-camel conversion knows no acronyms beyond "id",
// so anything with an acronym or a run-together word needs spelling out here.
// Tokens absent from this map use the default conversion.
var packageTypeKind = map[string]string{
	// Only the v2 registry format is supported, and only local repositories
	// carry the version in their Terraform name. Dropping it keeps the Kind
	// consistent with RemoteDockerRepository and VirtualDockerRepository.
	"docker_v2":        "Docker",
	"helmoci":          "HelmOCI",
	"pypi":             "PyPI",
	"terraformbackend": "TerraformBackend",
}

// kindFor returns the CRD Kind for a repository class and package type,
// e.g. ("local", "helmoci") -> "LocalHelmOCIRepository".
func kindFor(class, packageType string) string {
	pt, ok := packageTypeKind[packageType]
	if !ok {
		pt = name.NewFromSnake(packageType).Camel
	}
	return name.NewFromSnake(class).Camel + pt + "Repository"
}

// Matches reports whether the given Terraform resource name is a supported
// repository resource.
func Matches(name string) bool {
	return nameRe.MatchString(name)
}

// Configure registers a resource configurator for every supported local,
// remote and virtual repository resource.
//
// Upjet's default naming derives the API group from the second word of the
// Terraform name, which would scatter repositories across three groups
// (local, remote, virtual) with kinds like "NpmRepository". This provider
// instead exposes them all under a single "repository" group with the class
// folded into the kind, e.g. "RemoteDockerRepository".
func Configure(p *config.Provider) {
	for n := range p.Resources {
		m := nameRe.FindStringSubmatch(n)
		if m == nil {
			continue
		}
		class, packageType := m[1], m[2]
		kind := kindFor(class, packageType)
		isRemote := class == "remote"

		p.AddResourceConfigurator(n, func(r *config.Resource) {
			r.ShortGroup = "repository"
			r.Kind = kind

			if isRemote {
				// password_wo is a Terraform write-only attribute (provider
				// 12.11.8+), which requires Terraform 1.11 or later. This
				// provider is pinned to Terraform 1.5.7 because 1.6+ is BSL
				// licensed, so the write-only semantics are unavailable and the
				// attribute cannot behave correctly. Drop it rather than
				// exposing a field that silently misbehaves.
				delete(r.TerraformResource.Schema, "password_wo")
				delete(r.TerraformResource.Schema, "password_wo_version")
			}
		})
	}
}
