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

package storage

import (
	"context"

	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	storagev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/storage/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"

	awsstoragev1 "github.com/upbound/provider-aws/apis/s3/v1beta1"
	gcpstoragev1 "github.com/upbound/provider-gcp/apis/storage/v1beta1"
)

// StorageBucketsReconciler reconciles a StorageBuckets object
type StorageBucketsReconciler struct {
	utils.MissionClient
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets/finalizers,verbs=update

func (r *StorageBucketsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	bucket := &storagev1alpha1.StorageBuckets{}
	err := r.Get(ctx, req.NamespacedName, bucket)
	if err != nil {
		return ctrl.Result{}, err
	}

	mission, err := r.GetMission(ctx, bucket.Spec.MissionRef, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}
	result, err := r.ReconcileStorageBucket(ctx, bucket, &mission)
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *StorageBucketsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storagev1alpha1.StorageBuckets{}).
		Owns(&gcpstoragev1.Bucket{}).
		Owns(&awsstoragev1.Bucket{}).
		Complete(r)
}
