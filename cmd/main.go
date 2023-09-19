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

package main

import (
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	apischeme "k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	controllerscheme "sigs.k8s.io/controller-runtime/pkg/scheme"

	cpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	gcpcomputev1 "github.com/upbound/provider-gcp/apis/compute/v1beta1"
	gcpstoragev1 "github.com/upbound/provider-gcp/apis/storage/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"

	computev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/compute/v1alpha1"
	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	storagev1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/storage/v1alpha1"
	computecontroller "github.com/holy-tech/Mission-Control-Operator/internal/controller/compute"
	missioncontroler "github.com/holy-tech/Mission-Control-Operator/internal/controller/mission"
	storagecontroller "github.com/holy-tech/Mission-Control-Operator/internal/controller/storage"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(missionv1alpha1.AddToScheme(scheme))
	utilruntime.Must(computev1alpha1.AddToScheme(scheme))
	utilruntime.Must(storagev1alpha1.AddToScheme(scheme))

	crossplaneSchemeBuilder := &controllerscheme.Builder{GroupVersion: apischeme.GroupVersion{Group: "pkg.crossplane.io", Version: "v1"}}
	crossplaneSchemeBuilder.Register(
		&cpv1.Provider{},
		&cpv1.ProviderList{},
	)
	gcpSchemeBuilder := &controllerscheme.Builder{GroupVersion: apischeme.GroupVersion{Group: "gcp.upbound.io", Version: "v1beta1"}}
	gcpSchemeBuilder.Register(
		&gcpv1.ProviderConfig{},
		&gcpv1.ProviderConfigList{},
	)
	gcpComputeSchemeBuilder := &controllerscheme.Builder{GroupVersion: apischeme.GroupVersion{Group: "compute.gcp.upbound.io", Version: "v1beta1"}}
	gcpComputeSchemeBuilder.Register(
		&gcpcomputev1.Instance{},
		&gcpcomputev1.InstanceList{},
	)
	gcpStorageSchemeBuilder := &controllerscheme.Builder{GroupVersion: apischeme.GroupVersion{Group: "storage.gcp.upbound.io", Version: "v1beta1"}}
	gcpStorageSchemeBuilder.Register(
		&gcpstoragev1.Bucket{},
		&gcpstoragev1.BucketList{},
	)
	if err := crossplaneSchemeBuilder.AddToScheme(scheme); err != nil {
		os.Exit(1)
	}
	if err := gcpSchemeBuilder.AddToScheme(scheme); err != nil {
		os.Exit(1)
	}
	if err := gcpComputeSchemeBuilder.AddToScheme(scheme); err != nil {
		os.Exit(1)
	}
	if err := gcpStorageSchemeBuilder.AddToScheme(scheme); err != nil {
		os.Exit(1)
	}
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "28044579.mission-control.apis.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&missioncontroler.MissionReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("Mission"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Mission")
		os.Exit(1)
	}
	if err = (&missioncontroler.MissionKeyReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MissionKey")
		os.Exit(1)
	}
	if err = (&computecontroller.VirtualMachineReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "VirtualMachine")
		os.Exit(1)
	}
	if err = (&storagecontroller.StorageBucketsReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StorageBuckets")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
