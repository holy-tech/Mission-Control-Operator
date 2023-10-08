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

	runtime "k8s.io/apimachinery/pkg/runtime"
	record "k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	computev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/compute/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
	gcpcomputev1 "github.com/upbound/provider-gcp/apis/compute/v1beta1"
)

type VirtualMachineReconciler struct {
	utils.MissionClient
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
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

	mission, err := r.GetMission(ctx, vm.Spec.MissionRef.MissionName)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.ReconcileVirtualMachine(ctx, mission, vm)
	return ctrl.Result{}, err
}

func (r *VirtualMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1alpha1.VirtualMachine{}).
		Owns(&gcpcomputev1.Instance{}).
		Complete(r)
}
