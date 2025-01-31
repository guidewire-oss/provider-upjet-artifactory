// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	providerconfig "github.com/myorg/provider-jfrogartifactory/internal/controller/providerconfig"
	genericrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/repository/genericrepository"
	localnpmrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/repository/localnpmrepository"
	remotenpmrepository "github.com/myorg/provider-jfrogartifactory/internal/controller/repository/remotenpmrepository"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.Setup,
		genericrepository.Setup,
		localnpmrepository.Setup,
		remotenpmrepository.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
