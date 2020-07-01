package deployment

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"k8s.io/client-go/kubernetes"
)

type MissingDeployment struct {
	Name                           string
	ReadyReplicatError             error
	MinimumlReplicatAvailableError error
}

//HasDeploymentsInNamespace returns false and the list of deployment if some deployments are missing in the namespace.
//It returns an error if an error occurs while retreiving a deployment
//client: the client to use
//namespace: The namespace to search in
//expectedDeploymentNames: The deployment names to search
func HasDeploymentsInNamespace(client kubernetes.Interface, namespace string, expectedDeploymentNames []string) (has bool, missingDeployments []MissingDeployment, err error) {
	missingDeployments = make([]MissingDeployment, 0)
	has = true
	versionInfo, errDisco := client.Discovery().ServerVersion()
	if errDisco != nil {
		return false, missingDeployments, errDisco
	}
	klog.V(1).Infof("Server version info: %v", versionInfo)

	deployments := client.AppsV1().Deployments(namespace)

	for _, deploymentName := range expectedDeploymentNames {
		klog.V(1).Infof("Check if deployment %s exists", deploymentName)
		missingDeployment := MissingDeployment{
			Name: deploymentName,
		}
		missingDeploymentToBeAdded := false
		deployment, errGet := deployments.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if errGet != nil {
			if errors.IsNotFound(errGet) {
				missingDeploymentToBeAdded = true
				has = false
				missingDeployments = append(missingDeployments, missingDeployment)
				continue
			} else {
				klog.V(1).Infof("Error while retrieving deployment %s: %s", deploymentName, errGet.Error())
				return false, missingDeployments, errGet
			}
		}
		if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
			has = false
			missingDeploymentToBeAdded = true
			missingDeployment.ReadyReplicatError = fmt.Errorf("Expect %d for deployment %s but got %d Ready replicas",
				deployment.Status.Replicas,
				deploymentName,
				deployment.Status.ReadyReplicas)
			klog.Errorln(missingDeployment.ReadyReplicatError)
		}
		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" {
				if condition.Status != corev1.ConditionTrue {
					has = false
					missingDeploymentToBeAdded = true
					missingDeployment.MinimumlReplicatAvailableError = fmt.Errorf("Expect %s for deployment %s but got %s",
						condition.Status,
						deploymentName, corev1.ConditionTrue)
					klog.Errorln(missingDeployment.MinimumlReplicatAvailableError.Error())
				}
			}
		}
		if missingDeploymentToBeAdded {
			missingDeployments = append(missingDeployments, missingDeployment)
		}
	}

	return has, missingDeployments, err
}
