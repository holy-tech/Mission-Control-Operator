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
	"reflect"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	storagev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/storage/v1alpha1"

	gcpstoragev1 "github.com/upbound/provider-gcp/apis/storage/v1beta1"
)

// StorageBucketsReconciler reconciles a StorageBuckets object
type StorageBucketsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=storage.mission-control.apis.io,resources=storagebuckets/finalizers,verbs=update

func (r *StorageBucketsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	vm := &storagev1alpha1.StorageBuckets{}
	err := r.Get(ctx, req.NamespacedName, vm)
	if err != nil {
		return ctrl.Result{}, err
	}

	mission, err := r.GetMission(ctx, vm.Spec.MissionRef, req.Namespace)
	result, err := r.ReconcileStorageBucket(ctx, vm, &mission)
	return result, err
}

func (r *StorageBucketsReconciler) GetMission(ctx context.Context, missionName, missionNamespace string) (v1alpha1.Mission, error) {
	mission := v1alpha1.Mission{}
	err := r.Get(ctx, types.NamespacedName{Name: missionName, Namespace: missionNamespace}, &mission)
	return mission, err
}

func (r *StorageBucketsReconciler) ReconcileStorageBucket(ctx context.Context, bucket *storagev1alpha1.StorageBuckets, mission *v1alpha1.Mission) (ctrl.Result, error) {
	currentgcpbucket := gcpstoragev1.Bucket{}
	gcpbucket := gcpstoragev1.Bucket{
		ObjectMeta: v1.ObjectMeta{
			Name: bucket.Spec.ForProvider.Name,
		},
		Spec: gcpstoragev1.BucketSpec{
			ForProvider: gcpstoragev1.BucketParameters{
				Location: &bucket.Spec.ForProvider.Location,
			},
			ResourceSpec: cpcommonv1.ResourceSpec{
				ProviderConfigReference: &cpcommonv1.Reference{
					Name: "gcloud-provider",
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(bucket, &gcpbucket, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	err := r.Get(ctx, types.NamespacedName{Name: bucket.Spec.ForProvider.Name}, &currentgcpbucket)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, r.Create(ctx, &gcpbucket)
		}
		return ctrl.Result{}, err
	}
	if reflect.DeepEqual(currentgcpbucket.Spec, gcpbucket.Spec) {
		return ctrl.Result{}, nil
	}
	return reconcile.Result{}, r.Update(ctx, &currentgcpbucket)
}

// SetupWithManager sets up the controller with the Manager.
func (r *StorageBucketsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storagev1alpha1.StorageBuckets{}).
		Complete(r)
}
