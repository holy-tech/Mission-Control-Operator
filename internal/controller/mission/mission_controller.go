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
	"reflect"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	record "k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

type MissionReconciler struct {
	client.Client
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
	// Confirm that crossplane is installed in the kubernetes cluster
	if err := utils.ConfirmCRD(ctx, "providers.pkg.crossplane.io"); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", "Crossplane installation not found")
		return ctrl.Result{}, errors.New("could not find crossplane CRD \"Provider\"")
	}
	// Update mission providers installed status
	if err := ConfirmProviderConfigs(ctx, r, mission); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", err.Error())
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "Mission correctly connected to Crossplane")
	// Create ProviderConfig that resources will reference.
	if err := ReconcileProviderConfigs(ctx, r, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "ProviderConfig correctly created")
	// Confirm that mission key exists, if not create warning.
	if err := ConfirmMissionKeys(ctx, r, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "Mission keys correctly synced")
	return ctrl.Result{}, nil
}

func (r *MissionReconciler) GetProvider(ctx context.Context, providerName string) (*cpv1.Provider, error) {
	p := &cpv1.Provider{}
	err := r.Get(ctx, types.NamespacedName{Name: providerName}, p)
	return p, err
}

func (r *MissionReconciler) ApplyGenericProviderConfig(ctx context.Context, mission *missionv1alpha1.Mission, providerConfig, expectedProviderConfig utils.MissionObject) error {
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

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Owns(&gcpv1.ProviderConfig{}).
		Owns(&awsv1.ProviderConfig{}).
		Owns(&azrv1.ProviderConfig{}).
		Complete(r)
}
