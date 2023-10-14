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

	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

func ConfirmProvider(ctx context.Context, r *MissionReconciler, mission *missionv1alpha1.Mission) error {
	// Check that all the providers being used in specified mission
	// are installed in the cluster and are supported.
	for _, p := range mission.Spec.Packages {
		provider, err := GetProviderInstalled(ctx, r, mission, p.Provider)
		if err != nil {
			return err
		}
		UpdatePackageStatus(mission, provider)
		err = r.Status().Update(ctx, mission)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetProviderInstalled(ctx context.Context, r *MissionReconciler, mission *missionv1alpha1.Mission, providerName string) (*cpv1.Provider, error) {
	// Return provider after verifying that it is installed and supported by the software.
	// Returns error if provider is not installed or if not supported.
	if utils.Contains(utils.GetSupportedProviders(), providerName) {
		k8providerName := utils.ProviderMapping[providerName]
		p, err := r.GetProvider(ctx, k8providerName)
		if err != nil {
			message := fmt.Sprintf("Could not find provider %s, ensure provider is installed", k8providerName)
			return nil, errors.New(message)
		}
		return p, nil
	}
	message := fmt.Sprintf("Provider not allowed please choose of the following (%v)", utils.GetSupportedProviders())
	return nil, errors.New(message)
}
