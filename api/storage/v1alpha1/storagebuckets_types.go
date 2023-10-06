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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProviderData struct {
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

type StorageBucketMissionRef struct {
	MissionName string `json:"missionName,omitempty"`
	MissionKey  string `json:"missionKey,omitempty"`
}

type StorageBucketsSpec struct {
	MissionRef  StorageBucketMissionRef `json:"missionRef,omitempty"`
	ForProvider ProviderData            `json:"forProvider,omitempty"`
}

type StorageBucketsStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

type StorageBuckets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageBucketsSpec   `json:"spec,omitempty"`
	Status StorageBucketsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type StorageBucketsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageBuckets `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageBuckets{}, &StorageBucketsList{})
}
