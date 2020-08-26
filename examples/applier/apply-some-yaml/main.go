//Package main: This program shows how to create resources based on yamls template located in
//the different directory or bindata path or array of string.
package main

import (
	"flag"
	"os"

	// "github.com/open-cluster-management/library-go/examples/applier/bindata"
	"github.com/open-cluster-management/library-go/pkg/applier"
	libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func usage() {
	klog.Infof("Usage: apply-some-yaml -k kubeconfig\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	klog.InitFlags(nil)

	var kubeconfig = flag.String("k", "", "The path of the kubeconfig")
	var showHelp = flag.Bool("h", false, "Show help message")

	flag.Usage = usage
	flag.Parse()

	if *kubeconfig == "" {
		klog.Info("k is a mandatory argument")
		showUsageAndExit(0)
	}

	if *showHelp {
		showUsageAndExit(0)
	}

	err := applyYamlFile(*kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
}

func applyYamlFile(kubeconfig string) error {
	const directory = "../resources"
	//Create a reader on "../resources" directory
	klog.Infof("Creating the file reader %s", directory)
	yamlReader := applier.NewYamlFileReader(directory)
	//Other readers can be used
	//yamlReader := bindata.NewBindataReader()
	//yamlReader := applier.NewYamlStringReader(yamls,"---")

	//Create a templateProcessor with that reader
	klog.Infof("Creating TemplateProcessor...")
	tp, err := applier.NewTemplateProcessor(yamlReader, &applier.Options{})
	if err != nil {
		return err
	}
	//Create a client
	klog.Infof("Creating kubernetes client using kubeconfig located at %s", kubeconfig)
	client, err := libgoclient.NewDefaultClient(kubeconfig, client.Options{})
	if err != nil {
		return err
	}
	//Create an Applier
	klog.Info("Creating applier")
	applierOptions := &applier.ApplierOptions{
		Backoff: 	 &wait.Backoff{
				Steps:    5,
				Duration: 100 * time.Millisecond,
				Factor:   5.0,
				Jitter:   0.1,
			}
	}
	a, err := applier.NewApplier(tp, client, nil, nil, applier.DefaultKubernetesMerger, applierOptions)
	if err != nil {
		return err
	}
	//Defines the values
	values := struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}{
		ManagedClusterName:          "mycluster",
		ManagedClusterNamespace:     "mycluster",
		BootstrapServiceAccountName: "mybootstrapserviceaccount",
	}

	assetToBeApplied := []string{"yamlfilereader/namespace.yaml",
		"yamlfilereader/serviceaccount.yaml",
	}
	//Just to display what will be applied
	klog.Infof("Resources to be created: %v", assetToBeApplied)

	//template the resources
	us, err := tp.TemplateResourcesUnstructured(
		assetToBeApplied,
		values)
	if err != nil {
		return err
	}

	//Create the resources starting with path "yamlfilereader" in the reader
	//excluding "clusterrolebinding.yaml"
	//in a non-recursive way
	//and passing the values to replace
	klog.Info("Create or update resources")
	err = a.CreateOrUpdates(us)
	if err != nil {
		return err
	}
	klog.Infof("Resource deployed")
	return nil
}
