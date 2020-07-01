package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterregistryv1alpha1 "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
)

func TestAddFinalizer(t *testing.T) {
	testCluster := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	testCluster1 := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
				"test-finalizer",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	ExpectedtestCluster := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
				"test-finalizer",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	tests := []struct {
		name      string
		cluster   *clusterregistryv1alpha1.Cluster
		finalizer string
		Expected  *clusterregistryv1alpha1.Cluster
	}{
		{"add", testCluster, "test-finalizer", ExpectedtestCluster},
		{"don't add", testCluster1, "test-finalizer", ExpectedtestCluster},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddFinalizer(tt.cluster, tt.finalizer)
			assert.Equal(t, tt.cluster, tt.Expected)
		})
	}
}

func TestRemoveFinalizer(t *testing.T) {
	testCluster := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	testCluster1 := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
				"test-finalizer",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	ExpectedtestCluster := &clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"propagator.finalizer.mcm.ibm.com",
				"rcm-api.cluster",
			},
		},
		Status: clusterregistryv1alpha1.ClusterStatus{
			Conditions: []clusterregistryv1alpha1.ClusterCondition{
				{
					Status: corev1.ConditionTrue,
					Type:   clusterregistryv1alpha1.ClusterOK,
				},
			},
		},
	}
	tests := []struct {
		name      string
		cluster   *clusterregistryv1alpha1.Cluster
		finalizer string
		Expected  *clusterregistryv1alpha1.Cluster
	}{
		{"don't remove", testCluster, "test-finalizer", ExpectedtestCluster},
		{"remove", testCluster1, "test-finalizer", ExpectedtestCluster},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveFinalizer(tt.cluster, tt.finalizer)
			assert.Equal(t, tt.cluster, tt.Expected)
		})
	}
}
