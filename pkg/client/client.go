// Copyright (c) 2020 Red Hat, Inc.
package client

import (
	"github.com/open-cluster-management/library-go/pkg/config"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/klog"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

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
