package webhook_test

import (
	"os"
	"testing"

	k8scertutil "k8s.io/client-go/util/cert"

	webhook "github.com/open-cluster-management/library-go/pkg/webhook"
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
