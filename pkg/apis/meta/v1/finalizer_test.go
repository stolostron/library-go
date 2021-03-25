// Copyright Contributors to the Open Cluster Management project

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddFinalizer(t *testing.T) {
	testSecret := &corev1.Secret{
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
	}
	testSecret1 := &corev1.Secret{
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
	}
	ExpectedtestSecret := &corev1.Secret{
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
	}
	tests := []struct {
		name      string
		cluster   *corev1.Secret
		finalizer string
		Expected  *corev1.Secret
	}{
		{"add", testSecret, "test-finalizer", ExpectedtestSecret},
		{"don't add", testSecret1, "test-finalizer", ExpectedtestSecret},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddFinalizer(tt.cluster, tt.finalizer)
			assert.Equal(t, tt.cluster, tt.Expected)
		})
	}
}

func TestRemoveFinalizer(t *testing.T) {
	testSecret := &corev1.Secret{
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
	}
	testSecret1 := &corev1.Secret{
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
	}
	ExpectedtestSecret := &corev1.Secret{
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
	}
	tests := []struct {
		name      string
		cluster   *corev1.Secret
		finalizer string
		Expected  *corev1.Secret
	}{
		{"don't remove", testSecret, "test-finalizer", ExpectedtestSecret},
		{"remove", testSecret1, "test-finalizer", ExpectedtestSecret},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveFinalizer(tt.cluster, tt.finalizer)
			assert.Equal(t, tt.cluster, tt.Expected)
		})
	}
}

func TestHasFinalizer(t *testing.T) {
	testSecretWithFinalizer := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
			Finalizers: []string{
				"rcm-api.cluster",
			},
		},
	}
	testSecretWithoutFinalizer := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "test-cluster",
		},
	}
	type args struct {
		o         metav1.Object
		finalizer string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with finalizer",
			args: args{
				o:         testSecretWithFinalizer,
				finalizer: "rcm-api.cluster",
			},
			want: true,
		},
		{
			name: "without finalizer",
			args: args{
				o:         testSecretWithoutFinalizer,
				finalizer: "myfinalizer.cluster",
			},
			want: false,
		},
		{
			name: "with wrong finalizer",
			args: args{
				o:         testSecretWithFinalizer,
				finalizer: "wrong",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasFinalizer(tt.args.o, tt.args.finalizer); got != tt.want {
				t.Errorf("HasFinalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
