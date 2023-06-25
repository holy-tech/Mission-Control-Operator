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
	"os"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/v1alpha1"
)

type MissionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type Config struct {
	crdNamespace   string
	kubeconfigPath string
}

//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mission.mission-control.apis.io,resources=missions/finalizers,verbs=update

func (r *MissionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	clientConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	clientset, _ := apiextensionsclientset.NewForConfig(clientConfig)
	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, "providers.pkg.crossplane.io", v1.GetOptions{})

	if err != nil {
		return ctrl.Result{}, errors.New("Could not find crossplane CRD \"Provider\"")
	}

	ctrl.Log.Info("All seems correct")

	return ctrl.Result{}, nil
}

func (r *MissionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&missionv1alpha1.Mission{}).
		Complete(r)
}
