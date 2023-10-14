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
	"os"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	record "k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
)

type MissionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var ProviderMapping = map[string]string{
	"gcp":   "provider-gcp-family",
	"aws":   "provider-aws-family",
	"azure": "",
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
	if err := ConfirmCRD(ctx, "providers.pkg.crossplane.io"); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", "Crossplane installation not found")
		return ctrl.Result{}, errors.New("could not find crossplane CRD \"Provider\"")
	}
	// Check that the providers being used in specified mission are installed in the cluster and are supported
	if err := r.ConfirmProvider(ctx, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "Mission correctly connected to Crossplane")
	// Create ProviderConfig that resources will reference.
	if err := r.ReconcileProviderConfigs(ctx, mission); err != nil {
		return ctrl.Result{}, err
	}
	r.Recorder.Event(mission, "Normal", "Success", "ProviderConfig correctly created")
	// Confirm that mission key exists, if not create warning.
	if err := r.ConfirmMissionKeys(ctx, mission); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func ConfirmCRD(ctx context.Context, crdNameVersion string) error {
	clientConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	clientset, _ := apiextensionsclientset.NewForConfig(clientConfig)
	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdNameVersion, v1.GetOptions{})
	return err
}

func (r *MissionReconciler) ReconcilePackageStatus(ctx context.Context, mission *missionv1alpha1.Mission, provider *cpv1.Provider) error {
	if mission.Status.PackageStatus == nil {
		mission.Status.PackageStatus = map[string]missionv1alpha1.MissionPackageStatus{}
	}
	ps := mission.Status.PackageStatus[provider.Name]
	for _, c := range provider.Status.Conditions {
		if c.Type == "Installed" {
			ps.Installed = string(c.Status)
		}
	}
	mission.Status.PackageStatus[provider.Name] = ps
	return r.Status().Update(ctx, mission)
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Owns(&gcpv1.ProviderConfig{}).
		Complete(r)
}
