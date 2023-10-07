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
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
)

func (r *MissionKeyReconciler) CreateServiceAccount(ctx context.Context, req ctrl.Request, key *missionv1alpha1.MissionKey) error {
	sa := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := controllerutil.SetControllerReference(key, &sa, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, req.NamespacedName, &sa); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &sa)
		}
		return err
	}
	return nil
}
