// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	artifactorygroup "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/group/artifactorygroup"
	providerconfig "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/providerconfig"
	genericrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/genericrepository"
	localmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localmavenrepository"
	localnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localnpmrepository"
	remotemavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotemavenrepository"
	remotenpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotenpmrepository"
	virtualmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualmavenrepository"
	virtualnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualnpmrepository"
	artifactoryuser "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/user/artifactoryuser"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		artifactorygroup.Setup,
		providerconfig.Setup,
		genericrepository.Setup,
		localmavenrepository.Setup,
		localnpmrepository.Setup,
		remotemavenrepository.Setup,
		remotenpmrepository.Setup,
		virtualmavenrepository.Setup,
		virtualnpmrepository.Setup,
		artifactoryuser.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
