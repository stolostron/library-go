package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/open-cluster-management/library-go/pkg/applier"
	libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"github.com/open-cluster-management/library-go/pkg/templateprocessor"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Values map[string]interface{}

type Option struct {
	inFile         string
	outFile        string
	directory      string
	valuesPath     string
	kubeconfigPath string
	dryRun         bool
	prefix         string
	delete         bool
	timeout        int
	force          bool
}

func main() {
	var o Option
	klog.InitFlags(nil)
	flag.StringVar(&o.inFile, "f", "", "The file to process")
	flag.StringVar(&o.outFile, "o", "",
		"Output file. If set nothing will be applied but a file will be generate "+
			"which you can apply later with 'kubectl <create|apply|delete> -f")
	flag.StringVar(&o.directory, "d", ".", "The directory containing the templates, default '.'")
	flag.StringVar(&o.valuesPath, "values", "values.yaml", "The directory containing the templates, default 'values.yaml'")
	flag.StringVar(&o.kubeconfigPath, "k", "", "The kubeconfig file")
	flag.BoolVar(&o.dryRun, "dry-run", false, "if set only the rendered yaml will be shown, default false")
	flag.StringVar(&o.prefix, "p", "", "The prefix to add to each value names, for example 'Values'")
	flag.BoolVar(&o.delete, "delete", false,
		"if set only the resource defined in the yamls will be deleted, default false")
	flag.IntVar(&o.timeout, "t", 5, "Timeout in second to apply one resource, default 5 sec")
	flag.BoolVar(&o.force, "force", false, "If set, the finalizers will be removed before delete")
	flag.Parse()

	err := checkOptions(&o)
	if err != nil {
		fmt.Printf("Incorrect arguments: %s\n", err)
		os.Exit(1)
	}

	err = apply(o)
	if err != nil {
		fmt.Printf("Failed to apply due to error: %s\n", err)
		os.Exit(1)
	}
	if o.dryRun {
		fmt.Println("Dryrun successfully executed")
	} else {
		if o.outFile != "" {
			fmt.Println("Successfully processed")
		} else {
			fmt.Println("Sccessfully applied")
		}
	}
}

func checkOptions(o *Option) error {
	klog.V(2).Infof("-f: %s", o.inFile)
	klog.V(2).Infof("-d: %s", o.directory)
	if o.inFile != "" {
		o.directory = ""
	}
	if o.valuesPath == "" {
		return fmt.Errorf("-values is must be provided")
	}
	if o.inFile != "" && o.directory != "" {
		return fmt.Errorf("-f and -d are incompatible")
	}
	if o.outFile != "" &&
		(o.dryRun || o.delete || o.force) {
		return fmt.Errorf("-o is not compatible with -d, delete or force")
	}
	return nil
}

func apply(o Option) error {
	b, err := ioutil.ReadFile(filepath.Clean(o.valuesPath))
	if err != nil {
		return err
	}

	valuesc := &Values{}
	err = yaml.Unmarshal(b, valuesc)
	if err != nil {
		return err
	}

	values := Values{}
	if o.prefix != "" {
		values[o.prefix] = *valuesc
	} else {
		values = *valuesc
	}

	klog.V(5).Infof("values:\n%v", values)

	if o.outFile != "" {
		if o.inFile != "" {
			b, err := ioutil.ReadFile(filepath.Clean(o.inFile))
			if err != nil {
				return err
			}
			out, err := templateprocessor.TemplateBytes(b, values)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filepath.Clean(o.outFile), out, 0600)
		} else {
			templateReader := templateprocessor.NewYamlFileReader(o.directory)
			templateProcessor, err := templateprocessor.NewTemplateProcessor(templateReader, &templateprocessor.Options{})
			if err != nil {
				return err
			}
			outV, err := templateProcessor.TemplateResourcesInPathYaml("", []string{}, true, values)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filepath.Clean(o.outFile), []byte(templateprocessor.ConvertArrayOfBytesToString(outV)), 0600)
		}
	}
	var templateReader templateprocessor.TemplateReader
	if o.inFile != "" {
		b, err := ioutil.ReadFile(filepath.Clean(o.inFile))
		if err != nil {
			return err
		}
		templateReader = templateprocessor.NewYamlStringReader(string(b), templateprocessor.KubernetesYamlsDelimiter)
	} else {
		templateReader = templateprocessor.NewYamlFileReader(o.directory)
	}
	client, err := libgoclient.NewDefaultClient(o.kubeconfigPath, crclient.Options{})
	if err != nil {
		return err
	}
	applierOptions := &applier.Options{
		Backoff: &wait.Backoff{
			Steps:    4,
			Duration: 500 * time.Millisecond,
			Factor:   5.0,
			Jitter:   0.1,
			Cap:      time.Duration(o.timeout) * time.Second,
		},
		DryRun:      o.dryRun,
		ForceDelete: o.force,
	}
	if o.dryRun {
		client = crclient.NewDryRunClient(client)
	}
	a, err := applier.NewApplier(templateReader,
		&templateprocessor.Options{},
		client,
		nil,
		nil,
		applier.DefaultKubernetesMerger,
		applierOptions)
	if err != nil {
		return err
	}
	if o.delete {
		err = a.DeleteInPath("", nil, true, values)
	} else {
		err = a.CreateOrUpdateInPath("", nil, true, values)
	}
	if err != nil {
		return err
	}
	return nil
}
