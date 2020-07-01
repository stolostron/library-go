package deployment

import (
	"testing"

	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclientapps "k8s.io/client-go/kubernetes/fake"
)

func TestHaveDeploymentsInNamespace(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
	}

	deploymentMinAvailable := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Reason: "MinimumReplicasAvailable", Status: corev1.ConditionTrue},
			},
		},
	}
	deploymentNoMinAvailable := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Reason: "MinimumReplicasAvailable", Status: corev1.ConditionFalse},
			},
		},
	}

	client := fakeclientapps.NewSimpleClientset(deployment)
	clientMinAvailable := fakeclientapps.NewSimpleClientset(deploymentMinAvailable)
	clientNoMinAvailable := fakeclientapps.NewSimpleClientset(deploymentNoMinAvailable)

	type args struct {
		client                  kubernetes.Interface
		namespace               string
		expectedDeploymentNames []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "deployment exists",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: false,
		},
		{
			name: "all deployment not present",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment", "notexists"},
			},
			wantErr: true,
		},
		{
			name: "Deployment no minimum available",
			args: args{
				client:                  clientNoMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: true,
		},
		{
			name: "Deployment minimum available",
			args: args{
				client:                  clientMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HaveDeploymentsInNamespace(tt.args.client, tt.args.namespace, tt.args.expectedDeploymentNames); (err != nil) != tt.wantErr {
				t.Errorf("HaveDeploymentsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
