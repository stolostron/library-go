package e2etest

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

//GetClusters returns all clusters with a given tag
func GetClusters(tag string, clusters []Cluster) []*Cluster {
	filteredClusters := make([]*Cluster, 0)
	for i, cluster := range clusters {
		if tag, ok := cluster.Tags[tag]; ok {
			if tag {
				filteredClusters = append(filteredClusters, &clusters[i])
			}
		}
	}
	return filteredClusters
}
