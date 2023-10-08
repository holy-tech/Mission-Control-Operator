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

	missionv1alpha1 "github.com/holy-tech/Mission-Control-Operator/api/mission/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var c client.Client

var _ = Describe("Mission controller", func() {
	Context("Creating Mission", func() {
		It("Should apply mission definition", func() {
			By("Creating new mission")
			ctx := context.Background()
			mission := &missionv1alpha1.Mission{}
			Expect(k8sClient.Create(ctx, mission)).Should(Succeed())
		})
	})
})
