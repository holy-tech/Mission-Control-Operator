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

package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	record "k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
	utils "github.com/holy-tech/Mission-Control-Operator/controllers/utils"
)

type MissionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var ProviderMapping = map[string]string{
	"GCP":   "provider-gcp",
	"AWS":   "",
	"AZURE": "",
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MissionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mission := &missionv1alpha1.Mission{}
	err := r.Get(ctx, req.NamespacedName, mission)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Confirm that crossplane is installed in the kubernetes cluster
	if _, err := r.ConfirmCRD(ctx, "providers.pkg.crossplane.io"); err != nil {
		r.Recorder.Event(mission, "Warning", "Failed", "Crossplane installation not found")
		return ctrl.Result{}, errors.New("could not find crossplane CRD \"Provider\"")
	}

	// Check that the providers being used in specified mission are installed in the cluster and are supported
	for _, p := range mission.Spec.Packages {
		if !utils.ContainsString(utils.GetSupportedProviders(), p.Provider) {
			message := fmt.Sprintf("Provider %s is not supported, please use one of %v", p.Provider, utils.GetSupportedProviders())
			err := errors.New(message)
			return ctrl.Result{}, err
		}
		err := r.ConfirmProvider(ctx, mission, p.Provider)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Recorder.Event(mission, "Normal", "Success", "Mission correctly connected to Crossplane")

	// Create ProviderConfig that resources will reference.
	err = r.ReconcileProviderConfig(ctx, mission)
	if err != nil {
		r.Recorder.Event(mission, "Warning", "ProviderConfig not created", "Could not correctly create ProviderConfig resource.")
		return ctrl.Result{}, err
	}

	r.Recorder.Event(mission, "Normal", "Success", "ProviderConfig created")

	return ctrl.Result{}, nil
}

func (r *MissionReconciler) ConfirmProvider(ctx context.Context, mission *missionv1alpha1.Mission, providerName string) error {
	if utils.Contains(utils.GetValues(ProviderMapping), providerName) {
		k8providerName := ProviderMapping[providerName]
		p, err := r.GetProvider(ctx, k8providerName)
		if err != nil {
			message := fmt.Sprintf("Could not find provider %s, ensure provider is installed", k8providerName)
			r.Recorder.Event(mission, "Warning", "Provider Not Installed", message)
			return errors.New(message)
		}
		err = r.ReconcilePackageStatus(ctx, mission, p)
		if err != nil {
			return err
		}
	} else {
		message := fmt.Sprintf("Provider not allowed please choose of the following (%v)", utils.GetValues(ProviderMapping))
		r.Recorder.Event(mission, "Warning", "Provider Not Known", message)
		return errors.New(message)
	}
	return nil
}

func (r *MissionReconciler) GetProvider(ctx context.Context, providerName string) (*cpv1.Provider, error) {
	p := &cpv1.Provider{}
	err := r.Get(ctx, types.NamespacedName{Name: providerName}, p)
	return p, err
}

func (r *MissionReconciler) ConfirmCRD(ctx context.Context, crdNameVersion string) (*apiextensionsv1.CustomResourceDefinition, error) {
	clientConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	clientset, _ := apiextensionsclientset.NewForConfig(clientConfig)
	crd, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdNameVersion, v1.GetOptions{})
	return crd, err
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

func (r *MissionReconciler) ReconcileProviderConfig(ctx context.Context, mission *missionv1alpha1.Mission) error {
	providerConfig := &gcpv1.ProviderConfig{
		ObjectMeta: v1.ObjectMeta{
			Name: "testing config",
		},
		Spec: gcpv1.ProviderConfigSpec{
			ProjectID: "testing id",
			Credentials: gcpv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: "testing key",
						SecretReference: xpv1.SecretReference{
							Name:      "testing name",
							Namespace: "testing ns",
						},
					},
				},
			},
		},
	}
	if err := r.Get(ctx, types.NamespacedName{Name: mission.Name + "-GCP"}, providerConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, providerConfig)
		}
	}
	return r.Update(ctx, providerConfig)
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Complete(r)
}
