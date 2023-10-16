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

	utils "github.com/holy-tech/Mission-Control-Operator/internal/controller/utils"
)

func (k *MissionKey) GCPVerify() bool {
	k.GenericVerify()
	return true
}

func (k *MissionKey) AWSVerify() bool {
	k.GenericVerify()
	return true
}

func (k *MissionKey) AzureVerify() bool {
	k.GenericVerify()
	return true
}

func (k *MissionKey) GenericVerify() error {
	if !utils.Contains(utils.GetSupportedProviders(), k.Spec.Type) {
		message := fmt.Sprintf("Key of provider type %s is not supported, please use one of %v", k.Spec.Type, utils.GetSupportedProviders())
		err := errors.New(message)
		return err
	}
	return nil
}
