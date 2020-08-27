package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/open-cluster-management/library-go/pkg/applier"
	libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"github.com/open-cluster-management/library-go/pkg/templateprocessor"
	"gopkg.in/yaml.v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Values map[string]interface{}

func main() {
	var dir string
	var valuesPath string
	var kubeconfigPath string
	var dryRun bool
	var prefix string
	var delete bool
	var timeout int
	var force bool
	klog.InitFlags(nil)
	flag.StringVar(&dir, "d", ".", "The directory containing the templates, default '.'")
	flag.StringVar(&valuesPath, "values", "values.yaml", "The directory containing the templates, default 'values.yaml'")
	flag.StringVar(&kubeconfigPath, "k", "", "The kubeconfig file")
	flag.BoolVar(&dryRun, "dry-run", false, "if set only the rendered yaml will be shown, default false")
	flag.StringVar(&prefix, "p", "", "The prefix to add to each value names, for example 'Values'")
	flag.BoolVar(&delete, "delete", false, "if set only the resource defined in the yamls will be deleted, default false")
	flag.IntVar(&timeout, "t", 5, "Timeout in second to apply one resource, default 1 sec")
	flag.BoolVar(&force, "force", false, "If set, the finalizers will be removed before delete")
	flag.Parse()

	err := apply(dir, valuesPath, kubeconfigPath, prefix, timeout, dryRun, delete, force)
	if err != nil {
		fmt.Printf("Failed to apply due to error: %s", err)
		os.Exit(1)
	}
	if dryRun {
		fmt.Println("Dryrun successfully executed")
	} else {
		fmt.Println("Successfully applied")
	}
}

func apply(dir, valuesPath, kubeconfigPath, prefix string, timeout int, dryRun, delete, force bool) error {
	b, err := ioutil.ReadFile(filepath.Clean(valuesPath))
	if err != nil {
		return err
	}

	valuesc := &Values{}
	err = yaml.Unmarshal(b, valuesc)
	if err != nil {
		return err
	}

	values := Values{}
	if prefix != "" {
		values[prefix] = *valuesc
	} else {
		values = *valuesc
	}

	klog.V(5).Infof("values:\n%v", values)

	client, err := libgoclient.NewDefaultClient(kubeconfigPath, crclient.Options{})
	if err != nil {
		return err
	}
	applierOptions := &applier.Options{
		Backoff: &wait.Backoff{
			Steps:    4,
			Duration: 500 * time.Millisecond,
			Factor:   5.0,
			Jitter:   0.1,
			Cap:      time.Duration(timeout) * time.Second,
		},
		DryRun:      dryRun,
		ForceDelete: force,
	}
	if dryRun {
		client = crclient.NewDryRunClient(client)
	}
	a, err := applier.NewApplier(templateprocessor.NewYamlFileReader(dir),
		&templateprocessor.Options{},
		client,
		nil,
		nil,
		applier.DefaultKubernetesMerger,
		applierOptions)
	if err != nil {
		return err
	}
	if delete {
		err = a.DeleteInPath("", nil, true, values)
	} else {
		err = a.CreateOrUpdateInPath("", nil, true, values)
	}
	if err != nil {
		return err
	}
	return nil
}
