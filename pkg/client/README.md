# client

This pacakge contains several helper methods

## New client

You can create several types of kubernetes clients based on the kubeconfig.
The kubeconfig path can be provided or it use the KUBECONFIG if the path is "".
If in a cluster, then the config is retrieved from the cluster.

## Helpers

This package contains helper methods to check if resources are present such as:
- HaveServerResources
- HaveCRDs
- HaveDeploymentsInNamespace
