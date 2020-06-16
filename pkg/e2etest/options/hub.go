package options

import (
	"os"
	"path/filepath"
)

//GetHubKubeConfig returns the hub kubeconfig path if exists and empty string if not
//The kubeconfig file for the hub is supposed to be in <configDir>/<scenario>/<hubName>/kubeconfig
func GetHubKubeConfig(configDir, scenario string) string {
	if configDir == "" || scenario == "" {
		return ""
	}
	kubeConfigFilePath := filepath.Join(configDir, scenario, "kubeconfig.yaml")
	if _, err := os.Stat(kubeConfigFilePath); os.IsNotExist(err) {
		kubeConfigFilePath = ""
	}
	return kubeConfigFilePath
}
