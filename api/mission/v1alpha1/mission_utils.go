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
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	awsv1 "github.com/upbound/provider-aws/apis/v1beta1"
	azrv1 "github.com/upbound/provider-azure/apis/v1beta1"
	gcpv1 "github.com/upbound/provider-gcp/apis/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Mission) Convert2GCP(pkg *PackageConfig) *gcpv1.ProviderConfig {
	providerName := m.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &gcpv1.ProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ProviderConfig",
			APIVersion: "gcp.upbound.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: gcpv1.ProviderConfigSpec{
			ProjectID: pkg.ProjectID,
			Credentials: gcpv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	return providerConfig
}

func (m *Mission) Convert2AWS(pkg *PackageConfig) *awsv1.ProviderConfig {
	providerName := m.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &awsv1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: awsv1.ProviderConfigSpec{
			Credentials: awsv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	return providerConfig
}

func (m *Mission) Convert2Azure(pkg *PackageConfig) *azrv1.ProviderConfig {
	providerName := m.Name + "-" + strings.ToLower(pkg.Provider)
	providerConfig := &azrv1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: providerName,
		},
		Spec: azrv1.ProviderConfigSpec{
			Credentials: azrv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: pkg.Credentials.Key,
						SecretReference: xpv1.SecretReference{
							Name:      pkg.Credentials.Name,
							Namespace: pkg.Credentials.Namespace,
						},
					},
				},
			},
		},
	}
	return providerConfig
}

func (m *Mission) GCPVerify(packageId int) error {
	pkg := m.Spec.Packages[packageId]
	if pkg.ProjectID == "" {
		return errors.New("Project Id not filled for GCP package.")
	}
	return m.GenericVerify()
}

func (m *Mission) AWSVerify() error {
	return m.GenericVerify()
}

func (m *Mission) AzureVerify() error {
	return m.GenericVerify()
}

func (m *Mission) GenericVerify() error {
	return nil
}
