package providers

import (
	"bytes"
	"fmt"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-cluster-management/library-go/pkg/client"
	"github.com/open-cluster-management/library-go/pkg/e2etest/options"
	"k8s.io/klog"
)

type InstallerConfigAWS struct {
	Name          string
	BaseDnsDomain string
	SSHKey        string
	Region        string
}

// function for filling out installconfig template, takes installConfig struct and returns string for secret creation
func GetInstallConfigAWS(instConfig InstallerConfigAWS) string {
	const configTemplate = `
apiVersion: v1
metadata:
  name: {{.Name}}
baseDomain: {{.BaseDnsDomain}}
controlPlane:
  hyperthreading: Enabled
  name: master
  replicas: 3
  platform:
    aws:
      rootVolume:
        iops: 4000
        size: 500
        type: io1
      type: m4.xlarge
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
      aws:
      rootVolume:
        iops: 2000
        size: 500
        type: io1
      type: m4.large
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  aws:
    region: {{.Region}}
pullSecret: ""
sshKey: {{.SSHKey}}
`
	stdoutBuffer := new(bytes.Buffer)
	t := template.Must(template.New("configTemplate").Parse(configTemplate))
	err := t.Execute(stdoutBuffer, instConfig)
	if err != nil {
		return err.Error()
	}
	return stdoutBuffer.String()
}
func HaveServerResources(c options.Cluster, kubeconfig string, expectedAPIGroups []string) error {
	clientAPIExtension := client.NewKubeClientAPIExtension(c.MasterURL, kubeconfig, c.KubeContext)
	clientDiscovery := clientAPIExtension.Discovery()
	for _, apiGroup := range expectedAPIGroups {
		klog.V(1).Infof("Check if %s exists", apiGroup)
		_, err := clientDiscovery.ServerResourcesForGroupVersion(apiGroup)
		if err != nil {
			klog.V(1).Infof("Error while retrieving server resource %s: %s", apiGroup, err.Error())
			return err
		}
	}
	return nil
}

func HaveCRDs(c options.Cluster, kubeconfig string, expectedCRDs []string) error {
	clientAPIExtension := client.NewKubeClientAPIExtension(c.MasterURL, kubeconfig, c.KubeContext)
	clientAPIExtensionV1beta1 := clientAPIExtension.ApiextensionsV1beta1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, err := clientAPIExtensionV1beta1.CustomResourceDefinitions().Get(crd, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving crd %s: %s", crd, err.Error())
			return err
		}
	}
	return nil
}

func HaveDeploymentsInNamespace(c options.Cluster, kubeconfig string, namespace string, expectedDeploymentNames []string) error {

	client := client.NewKubeClient(c.MasterURL, kubeconfig, c.KubeContext)
	versionInfo, err := client.Discovery().ServerVersion()
	if err != nil {
		return err
	}
	klog.V(1).Infof("Server version info: %v", versionInfo)

	deployments := client.AppsV1().Deployments(namespace)

	for _, deploymentName := range expectedDeploymentNames {
		klog.V(1).Infof("Check if deployment %s exists", deploymentName)
		deployment, err := deployments.Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving deployment %s: %s", deploymentName, err.Error())
			return err
		}
		if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
			err = fmt.Errorf("Expect %d but got %d Ready replicas", deployment.Status.Replicas, deployment.Status.ReadyReplicas)
			klog.Errorln(err)
			return err
		}
		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" {
				if condition.Status != corev1.ConditionTrue {
					err = fmt.Errorf("Expect %s but got %s", condition.Status, corev1.ConditionTrue)
					klog.Errorln(err)
					return err
				}
			}
		}
	}

	return nil
}
