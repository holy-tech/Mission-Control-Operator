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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"
)

func (r *MissionReconciler) ReconcileProviderConfig(ctx context.Context, pkg *missionv1alpha1.PackageConfig, mission *missionv1alpha1.Mission) error {
	if pkg.Provider == "GCP" {
		return r.ReconcileProviderConfigGCP(ctx, pkg, mission)
	}
	if pkg.Provider == "AWS" {
		return r.ReconcileProviderConfigAWS(ctx, pkg, mission)
	}
	if pkg.Provider == "AZURE" {
		return r.ReconcileProviderConfigAzure(ctx, pkg, mission)
	}
	message := fmt.Sprintf("Provider %s not known", pkg.Provider)
	return errors.New(message)
}

func (r *MissionReconciler) ReconcileProviderConfigGCP(ctx context.Context, pkg *missionv1alpha1.PackageConfig, mission *missionv1alpha1.Mission) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &gcpv1.ProviderConfig{}
	expectedProviderConfig := &gcpv1.ProviderConfig{
		ObjectMeta: v1.ObjectMeta{
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
	if err := r.Get(ctx, types.NamespacedName{Name: providerName}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig, expectedProviderConfig) {
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}

func (r *MissionReconciler) ReconcileProviderConfigAWS(ctx context.Context, pkg *missionv1alpha1.PackageConfig, mission *missionv1alpha1.Mission) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &awsv1.ProviderConfig{}
	expectedProviderConfig := &awsv1.ProviderConfig{
		ObjectMeta: v1.ObjectMeta{
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
	if err := r.Get(ctx, types.NamespacedName{Name: providerName}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig, expectedProviderConfig) {
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}

func (r *MissionReconciler) ReconcileProviderConfigAzure(ctx context.Context, pkg *missionv1alpha1.PackageConfig, mission *missionv1alpha1.Mission) error {
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &azrv1.ProviderConfig{}
	expectedProviderConfig := &azrv1.ProviderConfig{
		ObjectMeta: v1.ObjectMeta{
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
	if err := r.Get(ctx, types.NamespacedName{Name: providerName}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(providerConfig, expectedProviderConfig) {
		providerConfig.Spec = expectedProviderConfig.Spec
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}
