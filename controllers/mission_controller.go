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

package controllers

import (
	"context"

	v1ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
)

type MissionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/finalizers,verbs=update

func (r *MissionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	crd := &v1ext.CustomResourceDefinition{}
	err := r.Get(ctx, types.NamespacedName{Name: "providers.pkg.crossplane.io"}, crd)
	if err != nil {
		return ctrl.Result{}, err
	}

	ctrl.Log.Info("All seems correct")

	return ctrl.Result{}, nil
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Complete(r)
}
