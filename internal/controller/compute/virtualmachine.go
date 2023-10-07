/*
Copyright 2023.

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

package compute

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	computev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/compute/v1alpha1"
	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	awscomputev1 "github.com/upbound/provider-aws/apis/ec2/v1beta1"
	gcpcomputev1 "github.com/upbound/provider-gcp/apis/compute/v1beta1"
)

func (r *VirtualMachineReconciler) ReconcileVirtualMachine(ctx context.Context, mission *v1alpha1.Mission, vm *computev1alpha1.VirtualMachine) error {
	keyName := vm.Spec.MissionRef.MissionKey
	missionKey, err := r.GetMissionKey(ctx, mission, keyName)
	if err != nil {
		return err
	}
	err = r.ReconcileVirtualMachineByProvider(ctx, mission, missionKey, vm)
	if err != nil {
		r.Recorder.Event(mission, "Warning", "ProviderConfig not created", "Could not correctly create ProviderConfig resource.")
		return err
	}
	return nil
}

func (r *VirtualMachineReconciler) ReconcileVirtualMachineByProvider(ctx context.Context, mission *v1alpha1.Mission, missionKey *v1alpha1.MissionKey, vm *computev1alpha1.VirtualMachine) error {
	var err error
	pkg := mission.Spec.Packages[0]
	if pkg.Provider == "GCP" {
		err = r.GetVirtualMachineGCP(ctx, mission, vm)
	} else if pkg.Provider == "AWS" {
		err = r.GetVirtualMachineAWS(ctx, mission, vm)
	} else {
		message := fmt.Sprintf("Provider %s not known", pkg.Provider)
		err = errors.New(message)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *VirtualMachineReconciler) GetVirtualMachineGCP(ctx context.Context, mission *v1alpha1.Mission, vm *computev1alpha1.VirtualMachine) error {
	// Create virtual machine config
	currentgcpvm := gcpcomputev1.Instance{}
	gcpvm := gcpcomputev1.Instance{
		ObjectMeta: v1.ObjectMeta{
			Name: vm.Spec.ForProvider.Name,
		},
		Spec: gcpcomputev1.InstanceSpec{
			ForProvider: gcpcomputev1.InstanceParameters{
				Zone:        &vm.Spec.ForProvider.Zone,
				MachineType: &vm.Spec.ForProvider.MachineType,
				BootDisk: []gcpcomputev1.BootDiskParameters{{
					InitializeParams: []gcpcomputev1.InitializeParamsParameters{{
						Image: &vm.Spec.ForProvider.Image,
					}},
				}},
				NetworkInterface: []gcpcomputev1.NetworkInterfaceParameters{{
					Network: &vm.Spec.ForProvider.Network,
				}},
			},
			ResourceSpec: cpcommonv1.ResourceSpec{
				ProviderConfigReference: &cpcommonv1.Reference{
					Name: "gcloud-provider",
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(vm, &gcpvm, r.Scheme); err != nil {
		return err
	}
	err := r.Get(ctx, types.NamespacedName{Name: vm.Spec.ForProvider.Name}, &currentgcpvm)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &gcpvm)
		}
		return err
	}
	if reflect.DeepEqual(currentgcpvm.Spec, gcpvm.Spec) {
		return nil
	}
	return r.Update(ctx, &currentgcpvm)
}

func (r *VirtualMachineReconciler) GetVirtualMachineAWS(ctx context.Context, mission *v1alpha1.Mission, vm *computev1alpha1.VirtualMachine) error {
	// Create virtual machine config
	currentawsvm := awscomputev1.Instance{}
	awsvm := awscomputev1.Instance{
		ObjectMeta: v1.ObjectMeta{
			Name: vm.Spec.ForProvider.Name,
		},
		Spec: awscomputev1.InstanceSpec{
			ForProvider:  awscomputev1.InstanceParameters{},
			ResourceSpec: cpcommonv1.ResourceSpec{},
		},
	}
	if err := controllerutil.SetControllerReference(vm, &awsvm, r.Scheme); err != nil {
		return err
	}
	err := r.Get(ctx, types.NamespacedName{Name: vm.Spec.ForProvider.Name}, &currentawsvm)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &awsvm)
		}
		return err
	}
	if reflect.DeepEqual(currentawsvm.Spec, awsvm.Spec) {
		return nil
	}
	return r.Update(ctx, &currentawsvm)
}
