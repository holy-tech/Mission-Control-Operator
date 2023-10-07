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
	"fmt"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	types "k8s.io/apimachinery/pkg/types"
)

func (r *MissionReconciler) ConfirmMissionKeys(ctx context.Context, mission *missionv1alpha1.Mission) error {
	for _, pkg := range mission.Spec.Packages {
		key := &missionv1alpha1.MissionKey{}
		err := r.Get(ctx, types.NamespacedName{Name: pkg.Credentials.Name, Namespace: pkg.Credentials.Namespace}, key)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				r.Recorder.Event(mission, "Warning", "Error looking for MissionKey", "Unexpected error while looking for MissionKey.")
				return err
			}
			message := fmt.Sprintf("Provider %s: Please ensure that MissionKey \"%s\" exists in namespace \"%s\".", pkg.Provider, pkg.Credentials.Name, pkg.Credentials.Namespace)
			r.Recorder.Event(mission, "Warning", "MissionKey not found", message)
		} else {
			message := fmt.Sprintf("MissionKey \"%s\" correctly linked in Namespace \"%s\".", pkg.Credentials.Name, pkg.Credentials.Namespace)
			r.Recorder.Event(mission, "Normal", "Success", message)
		}
	}
	return nil
}
