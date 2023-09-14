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

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gcpcomputev1 "github.com/upbound/provider-gcp/apis/compute/v1beta1"

	computev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/compute/v1alpha1"
)

type VirtualMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=compute.mission-control.apis.io,resources=virtualmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=compute.mission-control.apis.io,resources=virtualmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=compute.mission-control.apis.io,resources=virtualmachines/finalizers,verbs=update

func (r *VirtualMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	vm := &computev1alpha1.VirtualMachine{}
	err := r.Get(ctx, req.NamespacedName, vm)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.ReconcileVirtualMachine(ctx, vm, req)
}

func (r *VirtualMachineReconciler) ReconcileVirtualMachine(ctx context.Context, vm *computev1alpha1.VirtualMachine, req ctrl.Request) error {
	// Create virtual machine config
	gcpvm := gcpcomputev1.Instance{
		ObjectMeta: v1.ObjectMeta{
			Name: vm.Spec.ForProvider.Name,
		},
		Spec: gcpcomputev1.InstanceSpec{
			ForProvider: gcpcomputev1.InstanceParameters{
				Hostname:    &vm.Spec.ForProvider.Name,
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
		},
	}
	if err := controllerutil.SetControllerReference(vm, &gcpvm, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, req.NamespacedName, &gcpvm); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &gcpvm)
		}
		return err
	}
	return r.Update(ctx, &gcpvm)
}

func (r *VirtualMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1alpha1.VirtualMachine{}).
		Owns(&gcpcomputev1.Instance{}).
		Complete(r)
}
