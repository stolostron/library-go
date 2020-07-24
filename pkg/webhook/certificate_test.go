// Copyright 2019 The Kubernetes Authors.
//
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
	"os"
	"testing"

	webhook "github.com/open-cluster-management/library-go/pkg/webhook"
	k8scertutil "k8s.io/client-go/util/cert"
)

func TestGenerateSignedWebhookCertificates(t *testing.T) {
	webhookServiceName := "test-webhook-svc"
	webhookServiceNamespace := "default"

	certDir := "/tmp/tmp-cert"

	defer func() {
		os.RemoveAll(certDir)
	}()

	ca, err := webhook.GenerateWebhookCerts(certDir, webhookServiceNamespace, webhookServiceName)
	if err != nil {
		t.Errorf("Generate signed certificate failed, %v", err)
	}

	if ca == nil {
		t.Errorf("Generate signed certificate failed")
	}

	canReadCertAndKey, err := k8scertutil.CanReadCertAndKey("/tmp/tmp-cert/tls.crt", "/tmp/tmp-cert/tls.key")
	if err != nil {
		t.Errorf("Generate signed certificate failed, %v", err)
	}

	if !canReadCertAndKey {
		t.Errorf("Generate signed certificate failed")
	}
}
