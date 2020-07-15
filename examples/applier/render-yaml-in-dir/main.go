//Package main: This program shows how to create resources based on yamls template located in
//the same directory or bindata or array of string.
package main

import (
	"flag"
	"log"
	"os"

	// "github.com/open-cluster-management/library-go/examples/applier/bindata"
	"github.com/open-cluster-management/library-go/pkg/applier"
	"k8s.io/klog"
)

func usage() {
	log.Printf("Usage: render-yaml-in-dir\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	klog.InitFlags(nil)

	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	err := renderYamlFile()
	if err != nil {
		log.Fatal(err)
	}
}

func renderYamlFile() error {
	const directory = "../resources"
	//Create a reader on "../../resources" directory
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

	//Just to display what will be applied
	assetToBeApplied, err := tp.AssetNamesInPath(
		"yamlfilereader",
		[]string{"yamlfilereader/clusterrolebinding.yaml"},
		false)
	if err != nil {
		return err
	}
	klog.Infof("Resources to be created: %v", assetToBeApplied)

	//Render the resources starting with path "yamlfilereader" in the reader
	//excluding "clusterrolebinding.yaml"
	//in a non-recursive way
	//and passing the values to replace
	//The output is sorted
	klog.Info("Render resources\n")
	out, err := tp.TemplateAssetsInPathYaml(
		"yamlfilereader",
		[]string{"yamlfilereader/clusterrolebinding.yaml"},
		false,
		values)
	if err != nil {
		return err
	}
	klog.Infof("Generated resources yamls\n%s", applier.ConvertArrayOfBytesToString(out))
	return nil
}
