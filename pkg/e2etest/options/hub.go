package options

import (
	"os"
	"path/filepath"
)

//GetHubKubeConfig returns the hub kubeconfig path if exists and empty string if not
//The kubeconfig file for the hub is supposed to be in <configDir>/<scenario>/<hubName>/kubeconfig
func GetHubKubeConfig(configDir, scenario, hubName string) string {
	kubeConfigFilePath := filepath.Join(configDir, scenario, hubName, "kubeconfig.yaml")
	if _, err := os.Stat(kubeConfigFilePath); os.IsNotExist(err) {
		kubeConfigFilePath = ""
	}
	return kubeConfigFilePath
}
