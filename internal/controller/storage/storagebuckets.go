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
	"errors"
	"fmt"
	"reflect"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	storagev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/storage/v1alpha1"

	awsstoragev1 "github.com/upbound/provider-aws/apis/s3/v1beta1"
	gcpstoragev1 "github.com/upbound/provider-gcp/apis/storage/v1beta1"
)

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
					Name: mission.GetName() + "-" + strings.ToLower("GCP"),
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

func (r *StorageBucketsReconciler) ReconcileStorageBucketByProvider(ctx context.Context, bucket *storagev1alpha1.StorageBuckets, mission *v1alpha1.Mission) error {
	var err error
	pkg := mission.Spec.Packages[0]
	if pkg.Provider == "GCP" {
		err = r.GetStorageBucketGCP(ctx, bucket, mission)
	} else if pkg.Provider == "AWS" {
		err = r.GetStorageBucketAWS(ctx, bucket, mission)
	} else {
		message := fmt.Sprintf("Provider %s not known", pkg.Provider)
		err = errors.New(message)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *StorageBucketsReconciler) GetStorageBucketGCP(ctx context.Context, bucket *storagev1alpha1.StorageBuckets, mission *v1alpha1.Mission) error {
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
					Name: mission.GetName() + "-" + strings.ToLower("GCP"),
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(bucket, &gcpbucket, r.Scheme); err != nil {
		return err
	}
	err := r.Get(ctx, types.NamespacedName{Name: bucket.Spec.ForProvider.Name}, &currentgcpbucket)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &gcpbucket)
		}
		return err
	}
	if reflect.DeepEqual(currentgcpbucket.Spec, gcpbucket.Spec) {
		return nil
	}
	return r.Update(ctx, &currentgcpbucket)
}

func (r *StorageBucketsReconciler) GetStorageBucketAWS(ctx context.Context, bucket *storagev1alpha1.StorageBuckets, mission *v1alpha1.Mission) error {
	currentawsbucket := awsstoragev1.Bucket{}
	awsbucket := awsstoragev1.Bucket{
		ObjectMeta: v1.ObjectMeta{
			Name: bucket.Spec.ForProvider.Name,
		},
		Spec: awsstoragev1.BucketSpec{
			ForProvider: awsstoragev1.BucketParameters{
				Region: &bucket.Spec.ForProvider.Location,
			},
			ResourceSpec: cpcommonv1.ResourceSpec{
				ProviderConfigReference: &cpcommonv1.Reference{
					Name: mission.GetName() + "-" + strings.ToLower("AWS"),
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(bucket, &awsbucket, r.Scheme); err != nil {
		return err
	}
	err := r.Get(ctx, types.NamespacedName{Name: bucket.Spec.ForProvider.Name}, &currentawsbucket)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, &awsbucket)
		}
		return err
	}
	if reflect.DeepEqual(currentawsbucket.Spec, awsbucket.Spec) {
		return nil
	}
	return r.Update(ctx, &currentawsbucket)
}
