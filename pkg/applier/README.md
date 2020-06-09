# Introduction

The file [templateProcessor](../pkg/templateProcessor) contains an number of methods allowing you to render template yamls. 
The resources are read by an Go object satisfying the [TemplateReader](./templateProcessor.go) reader.  
The reader is embedded in a templateProcessor.TemplateProcessor object
The resources are sorted in order to be applied in a kubernetes environment using a templateProcessor.Client


## Implementing a reader

A reader will read assets from a data source. You can find [testreader_test.go](./testreader_test.go) an example of a reader which reads the data from memory.

A bindata implementation can be found [bindata](https://github.com/open-cluster-management/rcm-controller/pkg/bindata/bindatareader.go)


## How to use

The template parameters are passed using a `struct{}`

```
	values := struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}{
		ManagedClusterName:          instance.Name,
		ManagedClusterNamespace:     instance.Name,
		BootstrapServiceAccountName: instance.Name + bootstrapServiceAccountNamePostfix,
	}

	tp, err := applier.NewTemplateProcessor(bindata.NewBindataReader(), nil)
	if err != nil {
		return reconcile.Result{}, err
	}

	a, err := applier.NewApplier(tp, r.client, instance, r.scheme, merger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = a.CreateOrUpdateInPath(
		"test",
		nil,
		false,
		values,
	)
```

## Methods

In [templateProcessor](../pkg/templateProcessor) there are methods which templates the yamls, return them as a list of yamls or list of `unstructured.Unstructured`.
There are also methods that sort these templated yamls depending of their `kind`. The order is defined in `kindOrder` variable.
A method `CreateOrUpdateInPath` creates or update all resources localted in a specific path.

### Example 1: Generate a templated yaml

```
	values := struct {
		ManagedClusterName      string
		ManagedClusterNamespace string
	}{
		ManagedClusterName:      saNsN.Name,
		ManagedClusterNamespace: saNsN.Namespace,
	}
	tp, err := NewTemplateProcessor(NewTestReader(), nil)
	if err != nil {
		return nil, err
	}
	result, err := tp.TemplateAsset("hub/managedcluster/manifests/managedcluster-service-account.yaml", values)
	if err != nil {
		return nil, err
	}
```
The result contains a `[]byte` representing the templated yaml with the provided config.

### Example 2: Generate a list of templated yaml

```
	values := struct {
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

	tp, err := NewTemplateProcessor(NewTestReader(), nil)
	if err != nil {
		return nil, err
	}

	nucleusYAMLs, err := tp.TemplateAssets([]string{
		"klusterlet/namespace.yaml",
		"klusterlet/image_pull_secret.yaml",
		"klusterlet/bootstrap_secret.yaml",
		"klusterlet/cluster_role.yaml",
		"klusterlet/cluster_role_binding.yaml",
		"klusterlet/service_account.yaml",
		"klusterlet/operator.yaml",
	}, values)

```
nucleusYamls contains a non-sorted `[][]bytes` each element is the templated yamls using the provided config.

### Example 3: Generate a sorted list of yamls based using all templates in a given directory

```
	values := struct {
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

	tp, err := NewTemplateProcessor(NewTestReader(), nil)
	if err != nil {
		return nil, nil, err
	}

	nucleusYAMLs, err := tp.TemplateAssetsInPathYaml(
		"klusterlet", nil, false, values)
	if err != nil {
		return nil, nil, err
	}
```
The nucleusYAMls contains a `[][]byte` which is sorted based on all yamls in the `resources/klusterlet` (non-recursive) using the provided config.

### Example 4: Retreive a list of yamls

```
	tp, err := NewTemplateProcessor(NewTestReader(), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	crds, err = tp.Assets("klusterlet/crds", nil, true)
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
	values := struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}{
		ManagedClusterName:          instance.Name,
		ManagedClusterNamespace:     instance.Name,
		BootstrapServiceAccountName: instance.Name + bootstrapServiceAccountNamePostfix,
	}

	tp, err := NewTemplateProcessor(NewTestReader(), nil)
	if err != nil {
		return nil, nil, err
	}

	a, err := tp.NewApplier(a, r.client, instance, r.scheme, merger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = a.CreateOrUpdateInPath(
		"hub/managedcluster/manifests",
		nil,
		false,
		values,
	)

	if err != nil {
		return reconcile.Result{}, err
	}
```

This will create (in a sorted way) or update all resources located in the `hub/managedcluster/manifests` directory (non-recursive) except `hub/managedcluster/manifests/managedcluster-service-account.yaml`. A Merger function is passed as parameter to defind if the update must occur or not and how to merge the current resource with the new resource.
