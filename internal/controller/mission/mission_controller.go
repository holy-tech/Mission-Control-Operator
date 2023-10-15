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

	runtime "k8s.io/apimachinery/pkg/runtime"
	record "k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"

	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	clients "github.com/holy-tech/Mission-Control-Operator/internal/controller/clients"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

type MissionReconciler struct {
	clients.MissionClient
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MissionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mission := &missionv1alpha1.Mission{}
	if err := r.Get(ctx, req.NamespacedName, mission); err != nil {
		return ctrl.Result{}, err
	}
	// Ensure crossplane is installed in the kubernetes cluster
	if err := utils.ConfirmCRD(ctx, "providers.pkg.crossplane.io"); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", "Crossplane installation not found")
		return ctrl.Result{}, errors.New("could not find crossplane CRD \"Provider\"")
	}
	// Ensure crossplane providers are installed in the kubernetes cluster
	if err := ConfirmProviderConfigs(ctx, mission); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", err.Error())
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "Mission correctly connected to Crossplane")
	// Create ProviderConfigs that resources will reference.
	if err := ReconcileProviderConfigs(ctx, r, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "ProviderConfig correctly created")
	// Warn if mission keys are not created.
	if err := ConfirmMissionKeys(ctx, r, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "Mission keys correctly synced")
	return ctrl.Result{}, nil
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Owns(&gcpv1.ProviderConfig{}).
		Owns(&awsv1.ProviderConfig{}).
		Owns(&azrv1.ProviderConfig{}).
		Complete(r)
}
