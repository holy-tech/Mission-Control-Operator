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

	v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/controllers/utils"
)

type MissionKeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missionkeys/finalizers,verbs=update

func (r *MissionKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	key := &missionv1alpha1.MissionKey{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, key)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if secret and service account still exists if not create.
	secret := v1.Secret{
		Data: map[string][]byte{"keyfile": key.Spec.Data},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	sa := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &secret); err != nil {
		if errors.IsNotFound(err) {
			err := r.Create(ctx, &secret)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, err
	}
	if err = r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &sa); err != nil {
		if errors.IsNotFound(err) {
			err := r.Create(ctx, &sa)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, err
	}

	keyFinalizer := key.Spec.Name
	if key.ObjectMeta.DeletionTimestamp.IsZero() {
		if !utils.ContainsString(key.ObjectMeta.Finalizers, keyFinalizer) {
			key.ObjectMeta.Finalizers = append(key.ObjectMeta.Finalizers, keyFinalizer)
			if err := r.Update(context.Background(), key); err != nil {
				return reconcile.Result{}, err
			}
		}
	}
	if key.ObjectMeta.DeletionTimestamp != nil {
		// Delete secret and SA and check for err HERE

		key.ObjectMeta.Finalizers = utils.RemoveString(key.ObjectMeta.Finalizers, keyFinalizer)
		if err := r.Update(context.Background(), key); err != nil {
			return reconcile.Result{}, err
		}
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
