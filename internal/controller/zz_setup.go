// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	artifactorygroup "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/artifactorygroup"
	artifactoryuser "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/artifactoryuser"
	genericrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/genericrepository"
	localmavenrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/localmavenrepository"
	localnpmrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/localnpmrepository"
	remotemavenrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/remotemavenrepository"
	remotenpmrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/remotenpmrepository"
	virtualmavenrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/virtualmavenrepository"
	virtualnpmrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/jfrogartifactory/virtualnpmrepository"
	providerconfig "github.com/myorg/provider-jfrogartifactory/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		artifactorygroup.Setup,
		artifactoryuser.Setup,
		genericrepository.Setup,
		localmavenrepository.Setup,
		localnpmrepository.Setup,
		remotemavenrepository.Setup,
		remotenpmrepository.Setup,
		virtualmavenrepository.Setup,
		virtualnpmrepository.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
