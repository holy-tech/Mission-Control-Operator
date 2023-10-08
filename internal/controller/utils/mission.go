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

package utils

import (
	"context"
	"errors"
	"fmt"

	types "k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
)

type MissionClient struct {
	client.Client
}

func (r *MissionClient) GetMission(ctx context.Context, missionName string) (*v1alpha1.Mission, error) {
	mission := v1alpha1.Mission{}
	err := r.Get(ctx, types.NamespacedName{Name: missionName}, &mission)
	return &mission, err
}

func (r *MissionClient) GetMissionKey(ctx context.Context, mission *v1alpha1.Mission, keyName string) (*v1alpha1.MissionKey, error) {
	for _, pkg := range mission.Spec.Packages {
		missionkey := v1alpha1.MissionKey{}
		if pkg.Credentials.Name != keyName {
			continue
		}
		err := r.Get(ctx, types.NamespacedName{Name: pkg.Credentials.Name}, &missionkey)
		return &missionkey, err
	}
	msg := fmt.Sprintf("No credentials %s", keyName)
	return &v1alpha1.MissionKey{}, errors.New(msg)
}
