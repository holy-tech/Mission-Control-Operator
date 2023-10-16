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

package missionkeycontroller

import (
	"context"

	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	record "k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	clients "github.com/holy-tech/Mission-Control-Operator/internal/controller/clients"
)

type MissionKeyReconciler struct {
	clients.MissionClient
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys/finalizers,verbs=update

func (r *MissionKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	key := &missionv1alpha1.MissionKey{}
	err := r.Get(ctx, req.NamespacedName, key)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Ensure MissionKey is correct
	if err := key.GenericVerify(); err != nil {
		r.Recorder.Event(key, "Warning", "Failed", err.Error())
	}
	// Reconcile provider credential secrets
	if err := r.ReconcileSecret(ctx, key); err != nil {
		return ctrl.Result{}, err
	}
	// Reconcile service account for key usage
	if err := r.ReconcileServiceAccount(ctx, key); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MissionKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.MissionKey{}).
		Owns(&v1.Secret{}).
		Owns(&v1.ServiceAccount{}).
		Complete(r)
}
