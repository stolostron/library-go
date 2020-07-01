package deployment

import (
	"fmt"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	fakeclientapps "k8s.io/client-go/kubernetes/fake"
)

func Test_HasDeploymentsInNamespace(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
	}

	deploymentNotReady := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
		Status: appsv1.DeploymentStatus{
			Replicas:      1,
			ReadyReplicas: 0,
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
	clientNotReady := fakeclientapps.NewSimpleClientset(deploymentNotReady)

	type args struct {
		client                  kubernetes.Interface
		namespace               string
		expectedDeploymentNames []string
	}
	tests := []struct {
		name                   string
		args                   args
		wantHas                bool
		wantMissingDeployments []MissingDeployment
		wantErr                bool
	}{
		{
			name: "deployment exists",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantHas:                true,
			wantMissingDeployments: []MissingDeployment{},
			wantErr:                false,
		},
		{
			name: "all deployment not present",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment", "notexists"},
			},
			wantHas: false,
			wantMissingDeployments: []MissingDeployment{
				{Name: "notexists"},
			},
			wantErr: false,
		},
		{
			name: "Deployment not ready",
			args: args{
				client:                  clientNotReady,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantHas: false,
			wantMissingDeployments: []MissingDeployment{
				{
					Name: deploymentNoMinAvailable.Name,
					ReadyReplicatError: fmt.Errorf("Expect %d for deployment %s but got %d Ready replicas",
						deploymentNotReady.Status.Replicas,
						deploymentNotReady.Name,
						deploymentNotReady.Status.ReadyReplicas),
				},
			},
			wantErr: false,
		},
		{
			name: "Deployment no minimum available",
			args: args{
				client:                  clientNoMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantHas: false,
			wantMissingDeployments: []MissingDeployment{
				{
					Name: deploymentNoMinAvailable.Name,
					MinimumlReplicatAvailableError: fmt.Errorf("Expect %s for deployment %s but got %s",
						corev1.ConditionFalse,
						deploymentNoMinAvailable.Name, corev1.ConditionTrue),
				},
			},
			wantErr: false,
		},
		{
			name: "Deployment minimum available",
			args: args{
				client:                  clientMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantHas:                true,
			wantMissingDeployments: []MissingDeployment{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHas, gotMissingDeployments, err := HasDeploymentsInNamespace(tt.args.client, tt.args.namespace, tt.args.expectedDeploymentNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasDeploymentsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHas != tt.wantHas {
				t.Errorf("HasDeploymentsInNamespace() gotHas = %v, want %v", gotHas, tt.wantHas)
			}
			if !reflect.DeepEqual(gotMissingDeployments, tt.wantMissingDeployments) {
				t.Errorf("HasDeploymentsInNamespace() gotMissingDeployments = %v, want %v", gotMissingDeployments, tt.wantMissingDeployments)
			}
		})
	}
}
