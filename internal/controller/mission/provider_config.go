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

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"
)

func ReconcileProviderConfigs(ctx context.Context, r *MissionReconciler, mission *missionv1alpha1.Mission) error {
	for i, pkg := range mission.Spec.Packages {
		err := ReconcileProviderConfigByProvider(ctx, r, mission, i)
		if err != nil {
			message := fmt.Sprintf("Could not correctly create ProviderConfig resource %s.", pkg.Provider)
			r.Recorder.Event(mission, "Warning", "ProviderConfig not created", message)
			return err
		}
	}
	return nil
}

func ReconcileProviderConfigByProvider(ctx context.Context, r *MissionReconciler, mission *missionv1alpha1.Mission, packageId int) error {
	var err error
	pkg := &mission.Spec.Packages[packageId]
	if pkg.Provider == "gcp" {
		err = mission.GCPVerify(packageId)
		if err != nil {
			return err
		}
		err = r.ApplyGenericProviderConfig(
			ctx,
			mission,
			&gcpv1.ProviderConfig{},
			mission.Convert2GCP(pkg),
		)
	} else if pkg.Provider == "aws" {
		err = mission.AWSVerify()
		if err != nil {
			return err
		}
		err = r.ApplyGenericProviderConfig(
			ctx,
			mission,
			&awsv1.ProviderConfig{},
			mission.Convert2AWS(pkg),
		)
	} else if pkg.Provider == "azure" {
		err = mission.AzureVerify()
		if err != nil {
			return err
		}
		err = r.ApplyGenericProviderConfig(
			ctx,
			mission,
			&azrv1.ProviderConfig{},
			mission.Convert2Azure(pkg),
		)
	} else {
		message := fmt.Sprintf("Provider %s not known", pkg.Provider)
		err = errors.New(message)
	}
	if err != nil {
		return err
	}
	return nil
}

func ConfirmProviderConfigs(ctx context.Context, mission *missionv1alpha1.Mission) error {
	// Check that all providers being used in specified mission
	// are installed in the cluster and are supported.
	for _, p := range mission.Spec.Packages {
		providerCRD := fmt.Sprintf("providerconfigs.%s.upbound.io", p.Provider)
		if err := utils.ConfirmCRD(ctx, providerCRD); err != nil {
			return err
		}
	}
	return nil
}
