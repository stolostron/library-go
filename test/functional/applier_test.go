// +build functional

package functional_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	libgocrdv1 "github.com/open-cluster-management/library-go/pkg/apis/meta/v1/crd"
	"github.com/open-cluster-management/library-go/pkg/templateprocessor"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/open-cluster-management/library-go/pkg/applier"
	"gopkg.in/yaml.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Applier", func() {
	Context("Without Finalizer and no force", func() {
		It("Apply create/update resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample"), nil, clientHub, nil, nil, applier.DefaultKubernetesMerger, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.CreateOrUpdateInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrKlusterletAddonConfig).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).Should(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).Should(BeNil())

			Consistently(func() error {
				has, _, err := libgocrdv1.HasCRDs(clientAPIExt, []string{"klusterletaddonconfigs.agent.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).Should(BeNil())
		})

		It("Apply delete resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample"), nil, clientHub, nil, nil, applier.DefaultKubernetesMerger, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.DeleteInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrKlusterletAddonConfig).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				has, _, err := libgocrdv1.HasCRDs(clientAPIExt, []string{"klusterletaddonconfigs.agent.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).ShouldNot(BeNil())
		})
	})

	Context("With Finalizer and force", func() {
		It("Apply create/update resourcese", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample_with_finalizers"), nil, clientHub, nil, nil, applier.DefaultKubernetesMerger, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.CreateOrUpdateInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrKlusterletAddonConfig).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).Should(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).Should(BeNil())

			Consistently(func() error {
				has, _, err := libgocrdv1.HasCRDs(clientAPIExt, []string{"klusterletaddonconfigs.agent.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).Should(BeNil())
		})

		It("Apply delete resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample_with_finalizers"),
				nil,
				clientHub,
				nil,
				nil,
				applier.DefaultKubernetesMerger,
				&applier.Options{
					ForceDelete: true,
				})
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.DeleteInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrKlusterletAddonConfig).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				has, _, err := libgocrdv1.HasCRDs(clientAPIExt, []string{"klusterletaddonconfigs.agent.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).ShouldNot(BeNil())
		})
	})

})
