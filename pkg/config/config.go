// Copyright (c) 2020 Red Hat, Inc.

package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func LoadConfig(url, kubeconfig, context string) (*rest.Config, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	klog.V(5).Infof("Kubeconfig path %s\n", kubeconfig)
	// If we have an explicit indication of where the kubernetes config lives, read that.
	if kubeconfig != "" {
		if context == "" {
			// klog.V(5).Infof("clientcmd.BuildConfigFromFlags with %s and %s", url, kubeconfig)
			// Retreive the config for the current context
			if url == "" {
				config, err := clientcmd.LoadFromFile(kubeconfig)
				if err != nil {
					return nil, err
				}
				return clientcmd.NewDefaultClientConfig(
					*config,
					&clientcmd.ConfigOverrides{}).ClientConfig()
			}
			return clientcmd.BuildConfigFromFlags(url, kubeconfig)
		} else {
			return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
				&clientcmd.ConfigOverrides{
					CurrentContext: context,
				}).ClientConfig()
		}
	}
	// If not, try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		// log.Print("incluster\n")
		return c, nil
	}
	// If no in-cluster config, try the default location in the user's home directory.
	if usr, err := user.Current(); err == nil {
		klog.V(5).Infof("clientcmd.BuildConfigFromFlags for url %s using %s\n", url, filepath.Join(usr.HomeDir, ".kube", "config"))
		if c, err := clientcmd.BuildConfigFromFlags(url, filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not create a valid kubeconfig")

}
