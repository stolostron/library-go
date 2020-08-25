package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/open-cluster-management/library-go/pkg/applier"
	libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"gopkg.in/yaml.v1"
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
	klog.InitFlags(nil)
	flag.StringVar(&dir, "d", ".", "The directory containing the templates, default '.'")
	flag.StringVar(&valuesPath, "values", "values.yaml", "The directory containing the templates, default 'values.yaml'")
	flag.StringVar(&kubeconfigPath, "k", "", "The kubeconfig file")
	flag.BoolVar(&dryRun, "dry-run", false, "if set only the rendered yaml will be shown, default false")
	flag.StringVar(&prefix, "p", "", "The prefix to add to each value names, for example 'Values'")
	flag.Parse()

	err := apply(dir, valuesPath, kubeconfigPath, prefix, dryRun)
	if err != nil {
		klog.Errorf("Failed to apply due to error: %s", err)
		os.Exit(1)
	}
	if dryRun {
		fmt.Println("Dryrun successfully executed")
	} else {
		fmt.Println("Successfully applied")
	}
}

func apply(dir, valuesPath, kubeconfigPath, prefix string, dryRun bool) error {
	b, err := ioutil.ReadFile(valuesPath)
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

	tp, err := applier.NewTemplateProcessor(applier.NewYamlFileReader(dir), &applier.Options{})
	if err != nil {
		return err
	}

	client, err := libgoclient.NewDefaultClient(kubeconfigPath, crclient.Options{})
	if err != nil {
		return err
	}
	if dryRun {
		bb, err := tp.TemplateResourcesInPathYaml("", nil, true, values)
		if err != nil {
			return err
		}
		fmt.Print(applier.ConvertArrayOfBytesToString(bb))
		client = crclient.NewDryRunClient(client)
	}

	a, err := applier.NewApplier(tp, client, nil, nil, applier.DefaultKubernetesMerger, nil)
	if err != nil {
		return err
	}
	err = a.CreateOrUpdateInPath("", nil, true, values)
	if err != nil {
		return err
	}
	return nil
}
