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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	ctrl "sigs.k8s.io/controller-runtime"

	mgr "sigs.k8s.io/controller-runtime/pkg/manager"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

const (
	StartTimeout = 60 // seconds
)

var testEnv *envtest.Environment
var k8sManager mgr.Manager
var k8sClient client.Client

var (
	stop = ctrl.SetupSignalHandler()
)

func TestWebhookWireUp(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"WebHook Test",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	By("bootstrapping test environment")

	t := true
	if os.Getenv("TEST_USE_EXISTING_CLUSTER") == "true" {
		testEnv = &envtest.Environment{
			UseExistingCluster: &t,
		}
	} else {
		customAPIServerFlags := []string{"--disable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount," +
			"TaintNodesByCondition,Priority,DefaultTolerationSeconds,DefaultStorageClass,StorageObjectInUseProtection," +
			"PersistentVolumeClaimResize,ResourceQuota",
		}

		apiServerFlags := append([]string(nil), envtest.DefaultKubeAPIServerFlags...)
		apiServerFlags = append(apiServerFlags, customAPIServerFlags...)

		testEnv = &envtest.Environment{
			KubeAPIServerFlags: apiServerFlags,
		}
	}

	var err error
	// be careful, if we use shorthand assignment, the the cCfg will be a local variable
	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	k8sManager, err = mgr.New(cfg, mgr.Options{
		MetricsBindAddress: "0",
		Port:               testEnv.WebhookInstallOptions.LocalServingPort,
		Host:               testEnv.WebhookInstallOptions.LocalServingHost,
		CertDir:            testEnv.WebhookInstallOptions.LocalServingCertDir,
	})

	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(testEnv.Config, client.Options{})
	Expect(err).NotTo(HaveOccurred())

	go func() {
		Expect(k8sManager.Start(stop)).Should(Succeed())
	}()

	close(done)
}, StartTimeout)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
	Expect(testEnv.Stop()).ToNot(HaveOccurred())
})
