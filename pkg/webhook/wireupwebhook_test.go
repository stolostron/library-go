// Copyright 2019 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package webhook_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	webhook "github.com/open-cluster-management/library-go/pkg/webhook"
)

var _ = Describe("test if wireUp can create CA, service and validatingWebhookConfigration", func() {
	Context("given a manager it create svc and validating webhook config", func() {
		It("should create a service and ValidatingWebhookConfiguration", func() {
			testNs := "default"
			os.Setenv("POD_NAMESPACE", testNs)
			os.Setenv("DEPLOYMENT_LABEL", testNs)

			wbhName := "test-wbh"
			setWbhName := func(w *webhook.WebHookWireUp) {
				w.WebhookName = wbhName
			}

			wireUp, err := webhook.NewWebHookWireUp(k8sManager, stop, setWbhName)
			Expect(err).NotTo(HaveOccurred())

			caCert, err := wireUp.Attach()
			Expect(err).NotTo(HaveOccurred())

			wireUp.WireUpWebhookSupplymentryResource(caCert,
				schema.GroupVersionKind{Group: "", Version: "v1", Kind: "channels"},
				[]admissionv1.OperationType{admissionv1.Create})

			// give some time to allow the service and validtionconfig to come
			// up
			time.Sleep(3 * time.Second)

			wbhSvc := &corev1.Service{}
			svcKey := wireUp.WebHookeSvcKey
			Expect(k8sClient.Get(context.TODO(), svcKey, wbhSvc)).Should(Succeed())
			defer func() {
				Expect(k8sClient.Delete(context.TODO(), wbhSvc)).Should(Succeed())
			}()

			wbhCfg := &admissionv1.ValidatingWebhookConfiguration{}
			cfgKey := types.NamespacedName{Name: webhook.GetValidatorName(wbhName)}
			Expect(k8sClient.Get(context.TODO(), cfgKey, wbhCfg)).Should(Succeed())

			defer func() {
				Expect(k8sClient.Delete(context.TODO(), wbhCfg)).Should(Succeed())
			}()
		})
	})
})
