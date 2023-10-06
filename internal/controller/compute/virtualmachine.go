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
	"reflect"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	computev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/compute/v1alpha1"
	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	gcpcomputev1 "github.com/upbound/provider-gcp/apis/compute/v1beta1"
)

func (r *VirtualMachineReconciler) ReconcileVirtualMachine(ctx context.Context, vm *computev1alpha1.VirtualMachine, mission *v1alpha1.Mission) (ctrl.Result, error) {
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
		return ctrl.Result{}, err
	}
	err := r.Get(ctx, types.NamespacedName{Name: vm.Spec.ForProvider.Name}, &currentgcpvm)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, r.Create(ctx, &gcpvm)
		}
		return ctrl.Result{}, err
	}
	if reflect.DeepEqual(currentgcpvm.Spec, gcpvm.Spec) {
		return ctrl.Result{}, nil
	}
	return reconcile.Result{}, r.Update(ctx, &currentgcpvm)
}
