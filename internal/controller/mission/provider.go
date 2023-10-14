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

	types "k8s.io/apimachinery/pkg/types"

	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

func (r *MissionReconciler) ConfirmProvider(ctx context.Context, mission *missionv1alpha1.Mission) error {
	for _, p := range mission.Spec.Packages {
		if !utils.Contains(utils.GetSupportedProviders(), p.Provider) {
			message := fmt.Sprintf("Provider %s is not supported, please use one of %v", p.Provider, utils.GetSupportedProviders())
			err := errors.New(message)
			return err
		}
		err := r.ConfirmProviderInstalled(ctx, mission, p.Provider)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MissionReconciler) ConfirmProviderInstalled(ctx context.Context, mission *missionv1alpha1.Mission, providerName string) error {
	if utils.Contains(utils.GetValues(ProviderMapping), providerName) {
		k8providerName := ProviderMapping[providerName]
		p, err := r.GetProvider(ctx, k8providerName)
		if err != nil {
			message := fmt.Sprintf("Could not find provider %s, ensure provider is installed", k8providerName)
			return errors.New(message)
		}
		UpdatePackageStatus(mission, p)
		err = r.Status().Update(ctx, mission)
		if err != nil {
			return err
		}
	} else {
		message := fmt.Sprintf("Provider not allowed please choose of the following (%v)", utils.GetValues(ProviderMapping))
		return errors.New(message)
	}
	return nil
}

func (r *MissionReconciler) GetProvider(ctx context.Context, providerName string) (*cpv1.Provider, error) {
	p := &cpv1.Provider{}
	err := r.Get(ctx, types.NamespacedName{Name: providerName}, p)
	return p, err
}
