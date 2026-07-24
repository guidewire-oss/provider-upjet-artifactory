// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	artifactorygroup "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/artifactorygroup"
	artifactoryuser "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/artifactoryuser"
	genericrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/genericrepository"
	localmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/localmavenrepository"
	localnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/localnpmrepository"
	remotemavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/remotemavenrepository"
	remotenpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/remotenpmrepository"
	virtualmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/virtualmavenrepository"
	virtualnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/jfrogartifactory/virtualnpmrepository"
	providerconfig "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/cluster/providerconfig"
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

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		artifactorygroup.SetupGated,
		artifactoryuser.SetupGated,
		genericrepository.SetupGated,
		localmavenrepository.SetupGated,
		localnpmrepository.SetupGated,
		remotemavenrepository.SetupGated,
		remotenpmrepository.SetupGated,
		virtualmavenrepository.SetupGated,
		virtualnpmrepository.SetupGated,
		providerconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
