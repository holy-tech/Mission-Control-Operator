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

package missioncontroller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"
)

type ProviderConfigInterface interface {
	metav1.Object
	runtime.Object

	GetSpec()
	SetSpec()
}

func (r *MissionReconciler) ReconcileProviderConfigs(ctx context.Context, mission *missionv1alpha1.Mission) error {
	for i, pkg := range mission.Spec.Packages {
		err := r.ReconcileProviderConfigByProvider(ctx, mission, i, &pkg)
		if err != nil {
			r.Recorder.Event(mission, "Warning", "ProviderConfig not created", "Could not correctly create ProviderConfig resource.")
			return err
		}
	}
	return nil
}

func (r *MissionReconciler) ReconcileProviderConfigByProvider(ctx context.Context, mission *missionv1alpha1.Mission, packageId int, pkg *missionv1alpha1.PackageConfig) error {
	var err error
	if pkg.Provider == "gcp" {
		err = mission.GCPVerify(packageId)
		if err != nil {
			return err
		}
		err = r.GetProviderConfigGCP(ctx, mission, pkg)
	} else if pkg.Provider == "aws" {
		err = mission.AWSVerify()
		if err != nil {
			return err
		}
		err = r.GetProviderConfigAWS(ctx, mission, pkg)
	} else if pkg.Provider == "azure" {
		err = mission.AzureVerify()
		if err != nil {
			return err
		}
		err = r.GetProviderConfigAzure(ctx, mission, pkg)
	} else {
		message := fmt.Sprintf("Provider %s not known", pkg.Provider)
		err = errors.New(message)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *MissionReconciler) GetProviderConfigGCP(ctx context.Context, mission *missionv1alpha1.Mission, pkg *missionv1alpha1.PackageConfig) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &gcpv1.ProviderConfig{}
	expectedProviderConfig := &gcpv1.ProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ProviderConfig",
			APIVersion: "gcp.upbound.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: gcpv1.ProviderConfigSpec{
			ProjectID: pkg.ProjectID,
			Credentials: gcpv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(mission, expectedProviderConfig, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, types.NamespacedName{Name: expectedProviderConfig.GetName()}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig.Spec, expectedProviderConfig.Spec) {
		expectedProviderConfig.SetUID(providerConfig.GetUID())
		expectedProviderConfig.SetResourceVersion(providerConfig.GetResourceVersion())
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}

func (r *MissionReconciler) GetProviderConfigAWS(ctx context.Context, mission *missionv1alpha1.Mission, pkg *missionv1alpha1.PackageConfig) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &awsv1.ProviderConfig{}
	expectedProviderConfig := &awsv1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: awsv1.ProviderConfigSpec{
			Credentials: awsv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(mission, expectedProviderConfig, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, types.NamespacedName{Name: expectedProviderConfig.GetName()}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig.Spec, expectedProviderConfig.Spec) {
		expectedProviderConfig.SetUID(providerConfig.GetUID())
		expectedProviderConfig.SetResourceVersion(providerConfig.GetResourceVersion())
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}

func (r *MissionReconciler) GetProviderConfigAzure(ctx context.Context, mission *missionv1alpha1.Mission, pkg *missionv1alpha1.PackageConfig) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &azrv1.ProviderConfig{}
	expectedProviderConfig := &azrv1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: azrv1.ProviderConfigSpec{
			Credentials: azrv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(mission, expectedProviderConfig, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, types.NamespacedName{Name: expectedProviderConfig.GetName()}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig.Spec, expectedProviderConfig.Spec) {
		expectedProviderConfig.SetUID(providerConfig.GetUID())
		expectedProviderConfig.SetResourceVersion(providerConfig.GetResourceVersion())
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}
