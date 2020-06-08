# Introduction

The file [applier](../pkg/applier) contains an number of methods allowing you to render template yamls hold by a reader and also create these resources.

## Implementing a reader

A reader will read assets from a data source. You can find [testreader_test.go](./testreader_test.go) an example of a reader which reads the data from memory.

A bindata implementation can be found [bindata](https://github.com/open-cluster-management/rcm-controller/pkg/bindata)

## How to use

The template parameters are passed using a `struct{}`

```
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, err
	}
	config := struct {
		ManagedClusterName      string
		ManagedClusterNamespace string
	}{
		ManagedClusterName:      saNsN.Name,
		ManagedClusterNamespace: saNsN.Namespace,
	}
	result, err := a.TemplateAsset("hub/managedcluster/manifests/managedcluster-service-account.yaml")
	if err != nil {
		return nil, err
	}
```

## Methods

In [applier](../pkg/applier) there is methods which templates the yamls, return them as a list of yamls or list of `unstructured.Unstructured`.
There is also methods that sort these templated yamls depending of their `kind`. The order is defined in `kindOrder` variable.
A method `Create` creates all resources localted in a specific path

### Example 1: Generate a templated yaml

```
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, err
	}
	config := struct {
		ManagedClusterName      string
		ManagedClusterNamespace string
	}{
		ManagedClusterName:      saNsN.Name,
		ManagedClusterNamespace: saNsN.Namespace,
	}
	result, err := a.TemplateAsset("hub/managedcluster/manifests/managedcluster-service-account.yaml")
	if err != nil {
		return nil, err
	}
```
The result contains a `[]byte` representing the templated yaml with the provided config.

### Example 2: Generate a list of templated yaml

```
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, err
	}
	config := struct {
		KlusterletNamespace   string
		BootstrapSecretName   string
		BootstrapSecretToken  string
		BootstrapSecretCaCert string
		ImagePullSecretName   string
		ImagePullSecretData   string
		ImagePullSecretType   corev1.SecretType
	}{
		KlusterletNamespace:   klusterletNamespace,
		BootstrapSecretName:   managedCluster.Name,
		BootstrapSecretToken:  base64.StdEncoding.EncodeToString(bootStrapSecret.Data["token"]),
		BootstrapSecretCaCert: base64.StdEncoding.EncodeToString(bootStrapSecret.Data["ca.crt"]),
		ImagePullSecretName:   imagePullSecret.Name,
		ImagePullSecretData:   base64.StdEncoding.EncodeToString(imagePullSecret.Data[".dockerconfigjson"]),
		ImagePullSecretType:   imagePullSecret.Type,
	}

	nucleusYAMLs, err := a.TemplateAssets([]string{
		"klusterlet/namespace.yaml",
		"klusterlet/image_pull_secret.yaml",
		"klusterlet/bootstrap_secret.yaml",
		"klusterlet/cluster_role.yaml",
		"klusterlet/cluster_role_binding.yaml",
		"klusterlet/service_account.yaml",
		"klusterlet/operator.yaml",
	})

```
nucleusYamls contains a non-sorted `[][]bytes` each element is the templated yamls using the provided config.

### Example 3: Generate a sorted list of yamls based using all templates in a given directory

```
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	config := struct {
		KlusterletNamespace   string
		BootstrapSecretName   string
		BootstrapSecretToken  string
		BootstrapSecretCaCert string
		ImagePullSecretName   string
		ImagePullSecretData   string
		ImagePullSecretType   corev1.SecretType
	}{
		KlusterletNamespace:   klusterletNamespace,
		BootstrapSecretName:   managedCluster.Name,
		BootstrapSecretToken:  base64.StdEncoding.EncodeToString(bootStrapSecret.Data["token"]),
		BootstrapSecretCaCert: base64.StdEncoding.EncodeToString(bootStrapSecret.Data["ca.crt"]),
		ImagePullSecretName:   imagePullSecret.Name,
		ImagePullSecretData:   base64.StdEncoding.EncodeToString(imagePullSecret.Data[".dockerconfigjson"]),
		ImagePullSecretType:   imagePullSecret.Type,
	}

	nucleusYAMLs, err := a.TemplateAssetsInPathYaml(
		"klusterlet", nil, false)
	if err != nil {
		return nil, nil, err
	}
```
The nucleusYAMls contains a `[][]byte` which is sorted based on all yamls in the `resources/klusterlet` (non-recursive) using the provided config.

### Example 4: Retreive a list of yamls

```
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	crds, err = a.Assets("klusterlet/crds", nil, true)
	if err != nil {
		return nil, nil, err
	}
```
The crds contains a `[][]byte` (non-sorted) of all yamls found in `klusterlet/crds` directory and sub-directory using the provided config.

### Example 5: Create or update all resources defined in a directory

```
var merger bindata.Merger = func(current,
	new *unstructured.Unstructured,
) (
	future *unstructured.Unstructured,
	update bool,
) {
	if spec, ok := want.Object["spec"]; ok && 
	!reflect.DeepEqual(spec, current.Object["spec"]) {
		update = true
		current.Object["spec"] = spec
	}
	if rules, ok := want.Object["rules"]; ok && 
	!reflect.DeepEqual(rules, current.Object["rules"]) {
		update = true
		current.Object["rules"] = rules
	}
	if roleRef, ok := want.Object["roleRef"]; ok && 
	!reflect.DeepEqual(roleRef, current.Object["roleRef"]) {
		update = true
		current.Object["roleRef"] = roleRef
	}
	if subjects, ok := want.Object["subjects"]; ok && 
	!reflect.DeepEqual(subjects, current.Object["subjects"]) {
		update = true
		current.Object["subjects"] = subjects
	}
	return current, update
}
...
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		return nil, nil, err
	}

	c, err := applier.NewApplierClient(a, r.client, instance, r.scheme, merger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = c.CreateOrUpdateInPath(
		"hub/managedcluster/manifests",
		nil,
		false,
	)

	if err != nil {
		return reconcile.Result{}, err
	}
```

This will create (in a sorted way) or update all resources located in the `hub/managedcluster/manifests` directory (non-recursive) except `hub/managedcluster/manifests/managedcluster-service-account.yaml`. A Merger function is passed as parameter to defind if the update must occur or not and how to merge the current resource with the new resource.
