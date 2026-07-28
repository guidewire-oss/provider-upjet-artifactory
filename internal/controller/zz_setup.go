// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	artifactorygroup "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/group/artifactorygroup"
	artifactorypermissiontarget "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/permission/artifactorypermissiontarget"
	providerconfig "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/providerconfig"
	localdockerrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localdockerrepository"
	localgenericrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localgenericrepository"
	localgradlerepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localgradlerepository"
	localhelmocirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localhelmocirepository"
	localhelmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localhelmrepository"
	localmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localmavenrepository"
	localnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localnpmrepository"
	localpypirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localpypirepository"
	localterraformbackendrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localterraformbackendrepository"
	localterraformmodulerepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localterraformmodulerepository"
	localterraformproviderrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/localterraformproviderrepository"
	remotedockerrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotedockerrepository"
	remotegenericrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotegenericrepository"
	remotegradlerepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotegradlerepository"
	remotehelmocirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotehelmocirepository"
	remotehelmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotehelmrepository"
	remotemavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotemavenrepository"
	remotenpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotenpmrepository"
	remotepypirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remotepypirepository"
	remoteterraformrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/remoteterraformrepository"
	virtualdockerrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualdockerrepository"
	virtualgenericrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualgenericrepository"
	virtualgradlerepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualgradlerepository"
	virtualhelmocirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualhelmocirepository"
	virtualhelmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualhelmrepository"
	virtualmavenrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualmavenrepository"
	virtualnpmrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualnpmrepository"
	virtualpypirepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualpypirepository"
	virtualterraformrepository "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/repository/virtualterraformrepository"
	artifactoryuser "github.com/guidewire-oss/provider-jfrogartifactory/internal/controller/user/artifactoryuser"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		artifactorygroup.Setup,
		artifactorypermissiontarget.Setup,
		providerconfig.Setup,
		localdockerrepository.Setup,
		localgenericrepository.Setup,
		localgradlerepository.Setup,
		localhelmocirepository.Setup,
		localhelmrepository.Setup,
		localmavenrepository.Setup,
		localnpmrepository.Setup,
		localpypirepository.Setup,
		localterraformbackendrepository.Setup,
		localterraformmodulerepository.Setup,
		localterraformproviderrepository.Setup,
		remotedockerrepository.Setup,
		remotegenericrepository.Setup,
		remotegradlerepository.Setup,
		remotehelmocirepository.Setup,
		remotehelmrepository.Setup,
		remotemavenrepository.Setup,
		remotenpmrepository.Setup,
		remotepypirepository.Setup,
		remoteterraformrepository.Setup,
		virtualdockerrepository.Setup,
		virtualgenericrepository.Setup,
		virtualgradlerepository.Setup,
		virtualhelmocirepository.Setup,
		virtualhelmrepository.Setup,
		virtualmavenrepository.Setup,
		virtualnpmrepository.Setup,
		virtualpypirepository.Setup,
		virtualterraformrepository.Setup,
		artifactoryuser.Setup,
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
		artifactorypermissiontarget.SetupGated,
		providerconfig.SetupGated,
		localdockerrepository.SetupGated,
		localgenericrepository.SetupGated,
		localgradlerepository.SetupGated,
		localhelmocirepository.SetupGated,
		localhelmrepository.SetupGated,
		localmavenrepository.SetupGated,
		localnpmrepository.SetupGated,
		localpypirepository.SetupGated,
		localterraformbackendrepository.SetupGated,
		localterraformmodulerepository.SetupGated,
		localterraformproviderrepository.SetupGated,
		remotedockerrepository.SetupGated,
		remotegenericrepository.SetupGated,
		remotegradlerepository.SetupGated,
		remotehelmocirepository.SetupGated,
		remotehelmrepository.SetupGated,
		remotemavenrepository.SetupGated,
		remotenpmrepository.SetupGated,
		remotepypirepository.SetupGated,
		remoteterraformrepository.SetupGated,
		virtualdockerrepository.SetupGated,
		virtualgenericrepository.SetupGated,
		virtualgradlerepository.SetupGated,
		virtualhelmocirepository.SetupGated,
		virtualhelmrepository.SetupGated,
		virtualmavenrepository.SetupGated,
		virtualnpmrepository.SetupGated,
		virtualpypirepository.SetupGated,
		virtualterraformrepository.SetupGated,
		artifactoryuser.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
