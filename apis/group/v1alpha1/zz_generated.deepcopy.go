//go:build !ignore_autogenerated

// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroup) DeepCopyInto(out *ArtifactoryGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroup.
func (in *ArtifactoryGroup) DeepCopy() *ArtifactoryGroup {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArtifactoryGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupInitParameters) DeepCopyInto(out *ArtifactoryGroupInitParameters) {
	*out = *in
	if in.AdminPrivileges != nil {
		in, out := &in.AdminPrivileges, &out.AdminPrivileges
		*out = new(bool)
		**out = **in
	}
	if in.AutoJoin != nil {
		in, out := &in.AutoJoin, &out.AutoJoin
		*out = new(bool)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.DetachAllUsers != nil {
		in, out := &in.DetachAllUsers, &out.DetachAllUsers
		*out = new(bool)
		**out = **in
	}
	if in.ExternalID != nil {
		in, out := &in.ExternalID, &out.ExternalID
		*out = new(string)
		**out = **in
	}
	if in.PolicyManager != nil {
		in, out := &in.PolicyManager, &out.PolicyManager
		*out = new(bool)
		**out = **in
	}
	if in.Realm != nil {
		in, out := &in.Realm, &out.Realm
		*out = new(string)
		**out = **in
	}
	if in.RealmAttributes != nil {
		in, out := &in.RealmAttributes, &out.RealmAttributes
		*out = new(string)
		**out = **in
	}
	if in.ReportsManager != nil {
		in, out := &in.ReportsManager, &out.ReportsManager
		*out = new(bool)
		**out = **in
	}
	if in.UsersNames != nil {
		in, out := &in.UsersNames, &out.UsersNames
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.WatchManager != nil {
		in, out := &in.WatchManager, &out.WatchManager
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupInitParameters.
func (in *ArtifactoryGroupInitParameters) DeepCopy() *ArtifactoryGroupInitParameters {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupInitParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupList) DeepCopyInto(out *ArtifactoryGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ArtifactoryGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupList.
func (in *ArtifactoryGroupList) DeepCopy() *ArtifactoryGroupList {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArtifactoryGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupObservation) DeepCopyInto(out *ArtifactoryGroupObservation) {
	*out = *in
	if in.AdminPrivileges != nil {
		in, out := &in.AdminPrivileges, &out.AdminPrivileges
		*out = new(bool)
		**out = **in
	}
	if in.AutoJoin != nil {
		in, out := &in.AutoJoin, &out.AutoJoin
		*out = new(bool)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.DetachAllUsers != nil {
		in, out := &in.DetachAllUsers, &out.DetachAllUsers
		*out = new(bool)
		**out = **in
	}
	if in.ExternalID != nil {
		in, out := &in.ExternalID, &out.ExternalID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.PolicyManager != nil {
		in, out := &in.PolicyManager, &out.PolicyManager
		*out = new(bool)
		**out = **in
	}
	if in.Realm != nil {
		in, out := &in.Realm, &out.Realm
		*out = new(string)
		**out = **in
	}
	if in.RealmAttributes != nil {
		in, out := &in.RealmAttributes, &out.RealmAttributes
		*out = new(string)
		**out = **in
	}
	if in.ReportsManager != nil {
		in, out := &in.ReportsManager, &out.ReportsManager
		*out = new(bool)
		**out = **in
	}
	if in.UsersNames != nil {
		in, out := &in.UsersNames, &out.UsersNames
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.WatchManager != nil {
		in, out := &in.WatchManager, &out.WatchManager
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupObservation.
func (in *ArtifactoryGroupObservation) DeepCopy() *ArtifactoryGroupObservation {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupParameters) DeepCopyInto(out *ArtifactoryGroupParameters) {
	*out = *in
	if in.AdminPrivileges != nil {
		in, out := &in.AdminPrivileges, &out.AdminPrivileges
		*out = new(bool)
		**out = **in
	}
	if in.AutoJoin != nil {
		in, out := &in.AutoJoin, &out.AutoJoin
		*out = new(bool)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.DetachAllUsers != nil {
		in, out := &in.DetachAllUsers, &out.DetachAllUsers
		*out = new(bool)
		**out = **in
	}
	if in.ExternalID != nil {
		in, out := &in.ExternalID, &out.ExternalID
		*out = new(string)
		**out = **in
	}
	if in.PolicyManager != nil {
		in, out := &in.PolicyManager, &out.PolicyManager
		*out = new(bool)
		**out = **in
	}
	if in.Realm != nil {
		in, out := &in.Realm, &out.Realm
		*out = new(string)
		**out = **in
	}
	if in.RealmAttributes != nil {
		in, out := &in.RealmAttributes, &out.RealmAttributes
		*out = new(string)
		**out = **in
	}
	if in.ReportsManager != nil {
		in, out := &in.ReportsManager, &out.ReportsManager
		*out = new(bool)
		**out = **in
	}
	if in.UsersNames != nil {
		in, out := &in.UsersNames, &out.UsersNames
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.WatchManager != nil {
		in, out := &in.WatchManager, &out.WatchManager
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupParameters.
func (in *ArtifactoryGroupParameters) DeepCopy() *ArtifactoryGroupParameters {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupSpec) DeepCopyInto(out *ArtifactoryGroupSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
	in.InitProvider.DeepCopyInto(&out.InitProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupSpec.
func (in *ArtifactoryGroupSpec) DeepCopy() *ArtifactoryGroupSpec {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactoryGroupStatus) DeepCopyInto(out *ArtifactoryGroupStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactoryGroupStatus.
func (in *ArtifactoryGroupStatus) DeepCopy() *ArtifactoryGroupStatus {
	if in == nil {
		return nil
	}
	out := new(ArtifactoryGroupStatus)
	in.DeepCopyInto(out)
	return out
}
