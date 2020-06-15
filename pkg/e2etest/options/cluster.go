package options

import (
	"fmt"
	"os"
	"path/filepath"
)

//GetCluster returns the first cluster with a given tag
func GetCluster(tag string, clusters []Cluster) *Cluster {
	for _, cluster := range clusters {
		if tag, ok := cluster.Tags[tag]; ok {
			if tag {
				return &cluster
			}
		}
	}
	return nil
}

//GetKubeConfigs returns all for a given scenario kubeconfig.yaml path within a given configDir
func GetKubeConfigs(configDir string, scenario string) (map[string]string, error) {
	filteredKubeConfigs := make(map[string]string, 0)
	err := filepath.Walk(filepath.Join(configDir, scenario), func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if info.IsDir() {
				kubeConfigFilePath := filepath.Join(path + "kubeconfig.yaml")
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
