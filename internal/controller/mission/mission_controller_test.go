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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var c client.Client

var _ = Describe("GCP Mission controller", func() {
	Context("Creating GCP Mission", func() {
		It("Should apply gcp mission definition", func() {
			By("Creating new gcp mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mission-sample-gcp",
				},
				Spec: missionv1alpha1.MissionSpec{
					Packages: []missionv1alpha1.PackageConfig{{
						Provider:  "gcp",
						ProjectID: "made-up-project-id",
						Credentials: missionv1alpha1.CredentialConfig{
							Name:      "missionkey-sample-gcp",
							Namespace: "default",
							Key:       "creds",
						},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, mission)).Should(Succeed())
		})
	})
})

var _ = Describe("AWS Mission controller", func() {
	Context("Creating AWS Mission", func() {
		It("Should apply aws mission definition", func() {
			By("Creating new aws mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mission-sample-aws",
				},
				Spec: missionv1alpha1.MissionSpec{
					Packages: []missionv1alpha1.PackageConfig{{
						Provider: "aws",
						Credentials: missionv1alpha1.CredentialConfig{
							Name:      "missionkey-sample-aws",
							Namespace: "default",
							Key:       "creds",
						},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, mission)).Should(Succeed())
		})
	})
})

var _ = Describe("Azure Mission controller", func() {
	Context("Creating Azure Mission", func() {
		It("Should apply azure mission definition", func() {
			By("Creating new azure mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mission-sample-azure",
				},
				Spec: missionv1alpha1.MissionSpec{
					Packages: []missionv1alpha1.PackageConfig{{
						Provider: "azure",
						Credentials: missionv1alpha1.CredentialConfig{
							Name:      "missionkey-sample-azure",
							Namespace: "default",
							Key:       "creds",
						},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, mission)).Should(Succeed())
		})
	})
})
