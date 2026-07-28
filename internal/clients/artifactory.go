/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	"github.com/guidewire-oss/provider-jfrogartifactory/apis/v1beta1"
)

// Got the provider config fields from https://registry.terraform.io/providers/jfrog/artifactory/latest/docs
const (
	// KeyURL is the field containing the URL of Artifactory.
	KeyURL = "url"

	// KeyAccessToken is the field containing the Artifactory access token.
	KeyAccessToken = "access_token"

	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal jfrogartifactory credentials as JSON"
)

// TerraformSetupBuilder builds a terraform.SetupFn function which returns
// Terraform provider setup configuration. It resolves the referenced
// ProviderConfig for the managed resource, extracts credentials, and injects
// them into the Terraform provider configuration.
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		pcSpec, err := resolveProviderConfig(ctx, client, mg)
		if err != nil {
			return ps, errors.Wrap(err, "cannot resolve provider config")
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		if creds[KeyURL] == "" || creds[KeyAccessToken] == "" {
			println("missing required Artifactory credentials: url or access_token")
		}

		// Set credentials in Terraform provider configuration.
		ps.Configuration = map[string]any{
			KeyURL:         creds[KeyURL],
			KeyAccessToken: creds[KeyAccessToken],
		}
		return ps, nil
	}
}

// resolveProviderConfig resolves the ProviderConfig referenced by a
// cluster-scoped managed resource and returns its spec.
func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*v1beta1.ProviderConfigSpec, error) {
	managed, ok := mg.(resource.LegacyManaged) //nolint:staticcheck // cluster-scoped resources implement the legacy interface
	if !ok {
		return nil, errors.New("resource is not a cluster-scoped managed resource")
	}

	configRef := managed.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pc := &v1beta1.ProviderConfig{}
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(crClient, &v1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, managed); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return &pc.Spec, nil
}
