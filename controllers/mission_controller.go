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

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	record "k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	log "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
	"github.com/holy-tech/Mission-Control-Operator/controllers/utils"
)

type MissionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var Provider2CRD = map[string]string{
	"GCP":   "",
	"AWS":   "",
	"AZURE": "",
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MissionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	mission := &v1alpha1.Mission{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name}, mission)
	if err != nil {
		return ctrl.Result{}, err
	}

	if r.ConfirmCRD(ctx, "providers.pkg.crossplane.io") != nil {
		r.Recorder.Event(mission, "Warning", "Failed", "Crossplane installation not found")
		return ctrl.Result{}, errors.New("could not find crossplane CRD \"Provider\"")
	}

	for _, p := range mission.Spec.Packages {
		if utils.Contains(utils.GetValues(Provider2CRD), p) {
			providerCRD := Provider2CRD[p]
			if r.ConfirmCRD(ctx, providerCRD) != nil {
				message := fmt.Sprintf("Could not find provider %s, ensure provider is installed", p)
				r.Recorder.Event(mission, "Warning", "Provider Not Installed", message)
				return ctrl.Result{}, errors.New(message)
			}
		} else {
			message := fmt.Sprintf("Provider not allowed please choose of the following (%v)", utils.GetValues(Provider2CRD))
			r.Recorder.Event(mission, "Warning", "Provider Not Known", message)
			return ctrl.Result{}, errors.New(message)
		}
	}

	r.Recorder.Event(mission, "Normal", "Success", "Mission correctly connected to Crossplane")
	return ctrl.Result{}, nil
}

func (r *MissionReconciler) ConfirmCRD(ctx context.Context, crdNameVersion string) error {
	clientConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	clientset, _ := apiextensionsclientset.NewForConfig(clientConfig)
	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdNameVersion, v1.GetOptions{})
	return err
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Complete(r)
}
