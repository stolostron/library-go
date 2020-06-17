package helpers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/klog"
)

func HaveServerResources(client clientset.Interface, expectedAPIGroups []string) error {
	clientDiscovery := client.Discovery()
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

func HaveCRDs(client clientset.Interface, expectedCRDs []string) error {
	clientAPIExtensionV1beta1 := client.ApiextensionsV1beta1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, err := clientAPIExtensionV1beta1.CustomResourceDefinitions().Get(context.TODO(), crd, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving crd %s: %s", crd, err.Error())
			return err
		}
	}
	return nil
}

func HaveDeploymentsInNamespace(client kubernetes.Interface, namespace string, expectedDeploymentNames []string) error {
	versionInfo, err := client.Discovery().ServerVersion()
	if err != nil {
		return err
	}
	klog.V(1).Infof("Server version info: %v", versionInfo)

	deployments := client.AppsV1().Deployments(namespace)

	for _, deploymentName := range expectedDeploymentNames {
		klog.V(1).Infof("Check if deployment %s exists", deploymentName)
		deployment, err := deployments.Get(context.TODO(), deploymentName, metav1.GetOptions{})
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

//Apply a multi resources file to the cluster described by the url, kubeconfig and context.
//url of the cluster
//kubeconfig which contains the context
//context, the context to use
//yamlB, a byte array containing the resources file
// func Apply(url string, kubeconfig string, context string, yamlB []byte) error {
// 	yamls := strings.Split(string(yamlB), "---")
// 	clientKube := client.NewKubeClient(url, kubeconfig, context)
// 	clientAPIExtension := client.NewKubeClientAPIExtension(url, kubeconfig, context)
// // yamlFiles is an []string
// 	for _, f := range yamls {
// 		if len(strings.TrimSpace(f)) == 0 {
// 			continue
// 		}

// 		obj := &unstructured.Unstructured{}
// 		klog.V(5).Infof("obj:%v\n", obj.Object)
// 		err := yaml.Unmarshal([]byte(f), obj)
// 		if err != nil {
// 			return err
// 		}

// 		var kind string
// 		if v, ok := obj.Object["kind"]; !ok {
// 			return fmt.Errorf("kind attribute not found in %s", f)
// 		} else {
// 			kind = v.(string)
// 		}

// 		klog.V(5).Infof("kind: %s\n", kind)

// 		// now use switch over the type of the object
// 		// and match each type-case
// 		switch kind {
// 		case "CustomResourceDefinition":
// 			klog.V(5).Infof("Install CRD: %s\n", f)
// 			obj := &apiextensionsv1beta1.CustomResourceDefinition{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientAPIExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientAPIExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Create(obj)
// 			} else {
// 				existingObject.Spec = obj.Spec
// 				klog.Warningf("CRD %s already exists, updating!", existingObject.Name)
// 				_, err = clientAPIExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Update(existingObject)
// 			}
// 		case "Namespace":
// 			klog.V(5).Infof("Install %s: %s\n", kind, f)
// 			obj := &corev1.Namespace{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientKube.CoreV1().Namespaces().Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientKube.CoreV1().Namespaces().Create(obj)
// 			} else {
// 				obj.ObjectMeta = existingObject.ObjectMeta
// 				klog.Warningf("%s %s already exists, updating!", obj.Kind, obj.Name)
// 				_, err = clientKube.CoreV1().Namespaces().Update(existingObject)
// 			}
// 		case "ServiceAccount":
// 			klog.V(5).Infof("Install %s: %s\n", kind, f)
// 			obj := &corev1.ServiceAccount{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientKube.CoreV1().ServiceAccounts(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientKube.CoreV1().ServiceAccounts(obj.Namespace).Create(obj)
// 			} else {
// 				obj.ObjectMeta = existingObject.ObjectMeta
// 				klog.Warningf("%s %s/%s already exists, updating!", obj.Kind, obj.Namespace, obj.Name)
// 				_, err = clientKube.CoreV1().ServiceAccounts(obj.Namespace).Update(obj)
// 			}
// 		case "ClusterRoleBinding":
// 			klog.V(5).Infof("Install %s: %s\n", kind, f)
// 			obj := &rbacv1.ClusterRoleBinding{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientKube.RbacV1().ClusterRoleBindings().Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientKube.RbacV1().ClusterRoleBindings().Create(obj)
// 			} else {
// 				obj.ObjectMeta = existingObject.ObjectMeta
// 				klog.Warningf("%s %s/%s already exists, updating!", obj.Kind, obj.Namespace, obj.Name)
// 				_, err = clientKube.RbacV1().ClusterRoleBindings().Update(obj)
// 			}
// 		case "Secret":
// 			klog.V(5).Infof("Install %s: %s\n", kind, f)
// 			obj := &corev1.Secret{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientKube.CoreV1().Secrets(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientKube.CoreV1().Secrets(obj.Namespace).Create(obj)
// 			} else {
// 				obj.ObjectMeta = existingObject.ObjectMeta
// 				klog.Warningf("%s %s/%s already exists, updating!", obj.Kind, obj.Namespace, obj.Name)
// 				_, err = clientKube.CoreV1().Secrets(obj.Namespace).Update(obj)
// 			}
// 		case "Deployment":
// 			klog.V(5).Infof("Install %s: %s\n", kind, f)
// 			obj := &appsv1.Deployment{}
// 			err = yaml.Unmarshal([]byte(f), obj)
// 			if err != nil {
// 				return err
// 			}
// 			existingObject, errGet := clientKube.AppsV1().Deployments(obj.Namespace).Get(obj.Name, metav1.GetOptions{})
// 			if errGet != nil {
// 				_, err = clientKube.AppsV1().Deployments(obj.Namespace).Create(obj)
// 			} else {
// 				obj.ObjectMeta = existingObject.ObjectMeta
// 				klog.Warningf("%s %s/%s already exists, updating!", obj.Kind, obj.Namespace, obj.Name)
// 				_, err = clientKube.AppsV1().Deployments(obj.Namespace).Update(obj)
// 			}
// 		default:
// 			var resource string
// 			switch kind {
// 			case "Endpoint":
// 				klog.V(5).Infof("Install Endpoint: %s\n", f)
// 				resource = "endpoints"
// 			default:
// 				return fmt.Errorf("Resource %s not supported", kind)
// 			}
// 			var group string
// 			var version string
// 			if v, ok := obj.Object["apiVersion"]; !ok {
// 				return fmt.Errorf("apiVersion attribute not found in %s", f)
// 			} else {
// 				apiVersionArray := strings.Split(v.(string), "/")
// 				if len(apiVersionArray) != 2 {
// 					return fmt.Errorf("apiVersion malformed in %s", f)
// 				}
// 				group = apiVersionArray[0]
// 				version = apiVersionArray[1]
// 			}

// 			gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
// 			clientDynamic := client.NewKubeClientDynamic(url, kubeconfig, context)
// 			if ns := obj.GetNamespace(); ns != "" {
// 				existingObject, errGet := clientDynamic.Resource(gvr).Namespace(ns).Get(obj.GetName(), metav1.GetOptions{})
// 				if errGet != nil {
// 					_, err = clientDynamic.Resource(gvr).Namespace(ns).Create(obj, metav1.CreateOptions{})
// 				} else {
// 					obj.Object["metadata"] = existingObject.Object["metadata"]
// 					klog.Warningf("%s %s/%s already exists, updating!", obj.GetKind(), obj.GetNamespace(), obj.GetName())
// 					_, err = clientDynamic.Resource(gvr).Namespace(ns).Update(obj, metav1.UpdateOptions{})
// 				}
// 			} else {
// 				existingObject, errGet := clientDynamic.Resource(gvr).Get(obj.GetName(), metav1.GetOptions{})
// 				if errGet != nil {
// 					_, err = clientDynamic.Resource(gvr).Create(obj, metav1.CreateOptions{})
// 				} else {
// 					obj.Object["metadata"] = existingObject.Object["metadata"]
// 					klog.Warningf("%s %s already exists, updating!", obj.GetKind(), obj.GetName())
// 					_, err = clientDynamic.Resource(gvr).Update(obj, metav1.UpdateOptions{})
// 				}
// 			}
// 		}

// 		//&& !errors.IsAlreadyExists(err)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
