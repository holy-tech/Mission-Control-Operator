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

	types "k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
)

type MissionClient struct {
	client.Client
}

func (r *MissionClient) GetMission(ctx context.Context, missionName, missionNamespace string) (v1alpha1.Mission, error) {
	mission := v1alpha1.Mission{}
	err := r.Get(ctx, types.NamespacedName{Name: missionName, Namespace: missionNamespace}, &mission)
	return mission, err
}