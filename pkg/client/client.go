// Copyright (c) 2020 Red Hat, Inc.
package client

import (
	"context"
	"fmt"

	"github.com/open-cluster-management/library-go/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func NewDefaultClient(kubeconfig string, options client.Options) client.Client {
	return NewClient("", kubeconfig, "", options)
}

func NewClient(url, kubeconfig, context string, options client.Options) client.Client {
	klog.V(5).Infof("Create kubeclient for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	client, err := client.New(config, options)
	if err != nil {
		panic(err)
	}

	return client
}

func NewDefaultKubeClient(kubeconfig string) kubernetes.Interface {
	return NewKubeClient("", kubeconfig, "")
}

func NewKubeClient(url, kubeconfig, context string) kubernetes.Interface {
	klog.V(5).Infof("Create kubeclient for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func NewDefaultKubeClientDynamic(kubeconfig string) dynamic.Interface {
	return NewKubeClientDynamic("", kubeconfig, "")
}

func NewKubeClientDynamic(url, kubeconfig, context string) dynamic.Interface {
	klog.V(5).Infof("Create kubeclient dynamic for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func NewDefaultKubeClientAPIExtension(kubeconfig string) clientset.Interface {
	return NewKubeClientAPIExtension("", kubeconfig, "")
}

func NewKubeClientAPIExtension(url, kubeconfig, context string) clientset.Interface {
	klog.V(5).Infof("Create kubeclient apiextension for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func HaveServerResources(client clientset.Interface, expectedAPIGroups []string) error {
	clientDiscovery := client.Discovery()
	for _, apiGroup := range expectedAPIGroups {
		klog.V(1).Infof("Check if %s exists", apiGroup)
		_, err := clientDiscovery.ServerResourcesForGroupVersion(apiGroup)
		if err != nil {
			klog.V(1).Infof("Error while retrieving server resource %s: %s", apiGroup, err.Error())
			return err
		}
	}
	return nil
}

func HaveCRDs(client clientset.Interface, expectedCRDs []string) error {
	clientAPIExtensionV1beta1 := client.ApiextensionsV1beta1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, err := clientAPIExtensionV1beta1.CustomResourceDefinitions().Get(context.TODO(), crd, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving crd %s: %s", crd, err.Error())
			return err
		}
	}
	return nil
}

func HaveDeploymentsInNamespace(client kubernetes.Interface, namespace string, expectedDeploymentNames []string) error {
	versionInfo, err := client.Discovery().ServerVersion()
	if err != nil {
		return err
	}
	klog.V(1).Infof("Server version info: %v", versionInfo)

	deployments := client.AppsV1().Deployments(namespace)

	for _, deploymentName := range expectedDeploymentNames {
		klog.V(1).Infof("Check if deployment %s exists", deploymentName)
		deployment, err := deployments.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving deployment %s: %s", deploymentName, err.Error())
			return err
		}
		if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
			err = fmt.Errorf("Expect %d but got %d Ready replicas", deployment.Status.Replicas, deployment.Status.ReadyReplicas)
			klog.Errorln(err)
			return err
		}
		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" {
				if condition.Status != corev1.ConditionTrue {
					err = fmt.Errorf("Expect %s but got %s", condition.Status, corev1.ConditionTrue)
					klog.Errorln(err)
					return err
				}
			}
		}
	}

	return nil
}
