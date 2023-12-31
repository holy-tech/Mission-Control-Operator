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
)

var _ = Describe("Mission Creation controller", func() {
	Context("Creating Missions in different scenarios", func() {
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
		It("Should apply multi-provider mission definition", func() {
			By("Creating new multi-provider mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mission-sample-hybrid",
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
					}, {
						Provider: "aws",
						Credentials: missionv1alpha1.CredentialConfig{
							Name:      "missionkey-sample-aws",
							Namespace: "default",
							Key:       "creds",
						},
					}, {
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
		It("Should apply mission with multiple keys of one provider definition", func() {
			By("Creating new multiple key mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mission-sample-duplicate",
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
					}, {
						Provider:  "gcp",
						ProjectID: "made-up-project-id2",
						Credentials: missionv1alpha1.CredentialConfig{
							Name:      "missionkey-sample-gcp2",
							Namespace: "default",
							Key:       "creds",
						},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, mission)).Should(Succeed())
		})
		It("Should crash due to mission already existing", func() {
			By("Creating same mission")
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
			Expect(k8sClient.Create(ctx, mission)).ShouldNot(Succeed())
		})
	})
})
