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
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

func (k *MissionKey) Convert2Secret() *v1.Secret {
	return &v1.Secret{
		Data: map[string][]byte{"creds": k.Spec.Data},
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.GetName(),
			Namespace: k.GetNamespace(),
		},
	}
}

func (k *MissionKey) GCPVerify() error {
	return k.GenericVerify()
}

func (k *MissionKey) AWSVerify() error {
	return k.GenericVerify()
}

func (k *MissionKey) AzureVerify() error {
	return k.GenericVerify()
}

func (k *MissionKey) GenericVerify() error {
	if !utils.Contains(utils.GetSupportedProviders(), k.Spec.Type) {
		message := fmt.Sprintf("Key of provider type %s is not supported, please use one of %v", k.Spec.Type, utils.GetSupportedProviders())
		err := errors.New(message)
		return err
	}
	return nil
}
