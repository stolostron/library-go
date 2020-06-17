package options

import (
	"fmt"
	"os"
	"path/filepath"
)

//GetManagedClusterKubeConfigs returns all managedcluster kubeconfig for a given scenario
//The file path <configDir>/<scenario>/<clusterName>/kubeconfig.yaml are returned
func GetManagedClusterKubeConfigs(configDir, scenario string) (map[string]string, error) {
	if configDir == "" {
		return nil, fmt.Errorf("configDir not defined")
	}
	if scenario == "" {
		return nil, fmt.Errorf("scenario not defined")
	}
	scenarioPath := filepath.Join(configDir, scenario)
	if _, err := os.Stat(scenarioPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Missing scenario path %s", scenarioPath)
	}
	filteredKubeConfigs := make(map[string]string, 0)
	err := filepath.Walk(scenarioPath, func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if info.IsDir() && info.Name() != scenario {
				kubeConfigFilePath := filepath.Join(path, "kubeconfig.yaml")
				if _, err := os.Stat(kubeConfigFilePath); os.IsNotExist(err) {
					return fmt.Errorf("Missing file %s", kubeConfigFilePath)
				}
				filteredKubeConfigs[info.Name()] = kubeConfigFilePath
			}
		}
		return nil
	})
	return filteredKubeConfigs, err
}
