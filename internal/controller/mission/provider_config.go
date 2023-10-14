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
	types "k8s.io/apimachinery/pkg/types"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	"github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"
)

func (r *MissionReconciler) ReconcileProviderConfigs(ctx context.Context, mission *missionv1alpha1.Mission) error {
	for i, _ := range mission.Spec.Packages {
		err := r.ReconcileProviderConfigByProvider(ctx, mission, i)
		if err != nil {
			r.Recorder.Event(mission, "Warning", "ProviderConfig not created", "Could not correctly create ProviderConfig resource.")
			return err
		}
	}
	return nil
}

func (r *MissionReconciler) ReconcileProviderConfigByProvider(ctx context.Context, mission *missionv1alpha1.Mission, packageId int) error {
	var err error
	pkg := &mission.Spec.Packages[packageId]
	if pkg.Provider == "gcp" {
		err = mission.GCPVerify(packageId)
		if err != nil {
			return err
		}
		err = r.ApplyProviderConfigGCP(ctx, mission, packageId)
	} else if pkg.Provider == "aws" {
		err = mission.AWSVerify()
		if err != nil {
			return err
		}
		err = r.ApplyProviderConfigAWS(ctx, mission, packageId)
	} else if pkg.Provider == "azure" {
		err = mission.AzureVerify()
		if err != nil {
			return err
		}
		err = r.ApplyProviderConfigAzure(ctx, mission, packageId)
	} else {
		message := fmt.Sprintf("Provider %s not known", pkg.Provider)
		err = errors.New(message)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *MissionReconciler) ApplyProviderConfigGCP(ctx context.Context, mission *missionv1alpha1.Mission, packageId int) error {
	pkg := &mission.Spec.Packages[packageId]
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &gcpv1.ProviderConfig{}
	expectedProviderConfig := mission.Convert2GCP(providerName, pkg)
	return r.ApplyProviderConfig(ctx, mission, providerConfig, expectedProviderConfig)
}

func (r *MissionReconciler) ApplyProviderConfigAWS(ctx context.Context, mission *missionv1alpha1.Mission, packageId int) error {
	pkg := &mission.Spec.Packages[packageId]
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &awsv1.ProviderConfig{}
	expectedProviderConfig := mission.Convert2AWS(providerName, pkg)
	return r.ApplyProviderConfig(ctx, mission, providerConfig, expectedProviderConfig)
}

func (r *MissionReconciler) ApplyProviderConfigAzure(ctx context.Context, mission *missionv1alpha1.Mission, packageId int) error {
	pkg := &mission.Spec.Packages[packageId]
	providerName := mission.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &azrv1.ProviderConfig{}
	expectedProviderConfig := mission.Convert2Azure(providerName, pkg)
	return r.ApplyProviderConfig(ctx, mission, providerConfig, expectedProviderConfig)
}

func (r *MissionReconciler) ApplyProviderConfig(ctx context.Context, mission *missionv1alpha1.Mission, providerConfig, expectedProviderConfig utils.MissionObject) error {
	pcSpec := utils.GetValueOf(providerConfig, "Spec")
	epcSpec := utils.GetValueOf(expectedProviderConfig, "Spec")
	if pcSpec.Equal(reflect.Value{}) || epcSpec.Equal(reflect.Value{}) {
		return errors.New("Could not apply ProviderConfig")
	}
	if err := controllerutil.SetControllerReference(mission, expectedProviderConfig, r.Scheme); err != nil {
		return err
	}
	if err := r.Get(ctx, types.NamespacedName{Name: expectedProviderConfig.GetName()}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, expectedProviderConfig)
		}
	} else if !reflect.DeepEqual(pcSpec, epcSpec) {
		expectedProviderConfig.SetUID(providerConfig.GetUID())
		expectedProviderConfig.SetResourceVersion(providerConfig.GetResourceVersion())
		if err := utils.SetValueOf(providerConfig, "Spec", epcSpec); err != nil {
			return err
		}
		err := r.Update(ctx, providerConfig)
		return err
	}
	return nil
}
