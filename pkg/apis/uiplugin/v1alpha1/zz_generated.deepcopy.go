//go:build !ignore_autogenerated

/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertmanagerReference) DeepCopyInto(out *AlertmanagerReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertmanagerReference.
func (in *AlertmanagerReference) DeepCopy() *AlertmanagerReference {
	if in == nil {
		return nil
	}
	out := new(AlertmanagerReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentConfig) DeepCopyInto(out *DeploymentConfig) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentConfig.
func (in *DeploymentConfig) DeepCopy() *DeploymentConfig {
	if in == nil {
		return nil
	}
	out := new(DeploymentConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DistributedTracingConfig) DeepCopyInto(out *DistributedTracingConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DistributedTracingConfig.
func (in *DistributedTracingConfig) DeepCopy() *DistributedTracingConfig {
	if in == nil {
		return nil
	}
	out := new(DistributedTracingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingConfig) DeepCopyInto(out *LoggingConfig) {
	*out = *in
	out.LokiStack = in.LokiStack
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingConfig.
func (in *LoggingConfig) DeepCopy() *LoggingConfig {
	if in == nil {
		return nil
	}
	out := new(LoggingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStackReference) DeepCopyInto(out *LokiStackReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStackReference.
func (in *LokiStackReference) DeepCopy() *LokiStackReference {
	if in == nil {
		return nil
	}
	out := new(LokiStackReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringConfig) DeepCopyInto(out *MonitoringConfig) {
	*out = *in
	out.Alertmanager = in.Alertmanager
	out.ThanosQuerier = in.ThanosQuerier
	out.Perses = in.Perses
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringConfig.
func (in *MonitoringConfig) DeepCopy() *MonitoringConfig {
	if in == nil {
		return nil
	}
	out := new(MonitoringConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersesReference) DeepCopyInto(out *PersesReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersesReference.
func (in *PersesReference) DeepCopy() *PersesReference {
	if in == nil {
		return nil
	}
	out := new(PersesReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ThanosQuerierReference) DeepCopyInto(out *ThanosQuerierReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ThanosQuerierReference.
func (in *ThanosQuerierReference) DeepCopy() *ThanosQuerierReference {
	if in == nil {
		return nil
	}
	out := new(ThanosQuerierReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TroubleshootingPanelConfig) DeepCopyInto(out *TroubleshootingPanelConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TroubleshootingPanelConfig.
func (in *TroubleshootingPanelConfig) DeepCopy() *TroubleshootingPanelConfig {
	if in == nil {
		return nil
	}
	out := new(TroubleshootingPanelConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UIPlugin) DeepCopyInto(out *UIPlugin) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UIPlugin.
func (in *UIPlugin) DeepCopy() *UIPlugin {
	if in == nil {
		return nil
	}
	out := new(UIPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UIPlugin) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UIPluginList) DeepCopyInto(out *UIPluginList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]UIPlugin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UIPluginList.
func (in *UIPluginList) DeepCopy() *UIPluginList {
	if in == nil {
		return nil
	}
	out := new(UIPluginList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UIPluginList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UIPluginSpec) DeepCopyInto(out *UIPluginSpec) {
	*out = *in
	if in.Deployment != nil {
		in, out := &in.Deployment, &out.Deployment
		*out = new(DeploymentConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.TroubleshootingPanel != nil {
		in, out := &in.TroubleshootingPanel, &out.TroubleshootingPanel
		*out = new(TroubleshootingPanelConfig)
		**out = **in
	}
	if in.DistributedTracing != nil {
		in, out := &in.DistributedTracing, &out.DistributedTracing
		*out = new(DistributedTracingConfig)
		**out = **in
	}
	if in.Logging != nil {
		in, out := &in.Logging, &out.Logging
		*out = new(LoggingConfig)
		**out = **in
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UIPluginSpec.
func (in *UIPluginSpec) DeepCopy() *UIPluginSpec {
	if in == nil {
		return nil
	}
	out := new(UIPluginSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UIPluginStatus) DeepCopyInto(out *UIPluginStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UIPluginStatus.
func (in *UIPluginStatus) DeepCopy() *UIPluginStatus {
	if in == nil {
		return nil
	}
	out := new(UIPluginStatus)
	in.DeepCopyInto(out)
	return out
}
