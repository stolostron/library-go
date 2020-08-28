// +build functional

package functional_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Values map[string]interface{}

var (
	clientHub        client.Client
	clientHubDynamic dynamic.Interface
	clientAPIExt     clientset.Interface

	gvrKlusterletAddonConfig schema.GroupVersionResource
)

func init() {
	klog.SetOutput(GinkgoWriter)
	klog.InitFlags(nil)
}

var _ = BeforeSuite(func() {
	By("Setup Hub client")
	gvrKlusterletAddonConfig = schema.GroupVersionResource{Group: "agent.open-cluster-management.io", Version: "v1", Resource: "klusterletaddonconfigs"}

	var err error
	clientHub, err = libgoclient.NewDefaultClient("", client.Options{})
	Expect(err).To(BeNil())
	clientHubDynamic, err = libgoclient.NewDefaultKubeClientDynamic("")
	Expect(err).To(BeNil())
	clientAPIExt, err = libgoclient.NewDefaultKubeClientAPIExtension("")
	Expect(err).To(BeNil())
})

func TestFunctional(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional Suite")
}
