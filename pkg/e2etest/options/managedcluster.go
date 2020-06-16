package options

import (
	"fmt"
	"os"
	"path/filepath"
)

//GetClusterKubeConfigs returns all managedcluster kubeconfig for a given scenario, except for the hubName
//The file path <configDir>/<scenario>/<clusterName>/kubeconfig.yaml are returned
func GetClusterKubeConfigs(configDir, scenario, hubName string) (map[string]string, error) {
	filteredKubeConfigs := make(map[string]string, 0)
	err := filepath.Walk(filepath.Join(configDir, scenario), func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if info.IsDir() {
				if info.Name() != hubName {
					kubeConfigFilePath := filepath.Join(path + "kubeconfig.yaml")
					if _, err := os.Stat(kubeConfigFilePath); os.IsNotExist(err) {
						return fmt.Errorf("Missing file %s", kubeConfigFilePath)
					}
					filteredKubeConfigs[info.Name()] = kubeConfigFilePath
				}
			}
		}
		return nil
	})
	return filteredKubeConfigs, err
}
