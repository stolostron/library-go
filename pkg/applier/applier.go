package applier

import (
	"context"
	goerr "errors"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Applier structure to access kubernetes through the applier
type Applier struct {
	//An TemplateProcessor
	templateProcessor *TemplateProcessor
	//The client-go kubernetes client
	client client.Client
	//The owner of the created object
	owner metav1.Object
	//The scheme
	scheme *runtime.Scheme
	//A merger defining how two objects must be merged
	merger Merger
	//applier options for the applier
	applierOptions *ApplierOptions
}

//Options defines for the available options for the applier
type ApplierOptions struct {
	ClientCreateOption []client.CreateOption
	ClientUpdateOption []client.UpdateOption
	Backoff            *wait.Backoff
}

//NewApplier creates a new client to access kubernetes through the applier.
//applier: An applier
//client: The client-go client to use when applying the resources.
//owner: The object owner for the setControllerReference, the reference is not if nil.
//scheme: The object scheme for the setControllerReference, the reference is not if nil.
//merger: The function implementing the way how the resources must be merged
func NewApplier(
	templateProcessor *TemplateProcessor,
	client client.Client,
	owner metav1.Object,
	scheme *runtime.Scheme,
	merger Merger,
	applierOptions *ApplierOptions,
) (*Applier, error) {
	if templateProcessor == nil {
		return nil, goerr.New("applier is nil")
	}
	if client == nil {
		return nil, goerr.New("client is nil")
	}
	if applierOptions == nil {
		applierOptions = &ApplierOptions{}
	}
	if applierOptions.Backoff == nil {
		applierOptions.Backoff = &wait.Backoff{
			Steps:    5,
			Duration: 100 * time.Millisecond,
			Factor:   5.0,
			Jitter:   0.1,
		}
	}
	return &Applier{
		templateProcessor: templateProcessor,
		client:            client,
		owner:             owner,
		scheme:            scheme,
		merger:            merger,
		applierOptions:    applierOptions,
	}, nil
}

//Merger merges the `current` and the `want` resources into one resource which will be use for to update.
// If `update` is true than the update will be executed.
type Merger func(current,
	new *unstructured.Unstructured,
) (
	future *unstructured.Unstructured,
	update bool,
)

var rootAttributes = []string{
	"spec",
	"rules",
	"roleRef",
	"subjects",
	"secrets",
	"imagePullSecrets",
	"automountServiceAccountToken"}

//DefaultKubernetesMerger merges kubernetes runtime.Object
//It merges the spec, rules, roleRef, subjects root attribute of a runtime.Object
//For example a CLusterRoleBinding has a subjects and roleRef fields and so the old
//subjects and roleRef fields from the ClusterRoleBinding will be replaced by the new values.
var DefaultKubernetesMerger Merger = func(current,
	new *unstructured.Unstructured,
) (
	future *unstructured.Unstructured,
	update bool,
) {
	for _, r := range rootAttributes {
		if newValue, ok := new.Object[r]; ok {
			if !reflect.DeepEqual(newValue, current.Object[r]) {
				update = true
				current.Object[r] = newValue
			}
		} else {
			if _, ok := current.Object[r]; ok {
				current.Object[r] = nil
			}
		}
	}
	return current, update
}

//CreateOrUpdateInPath creates or updates the assets found in the path and
// subpath if recursive is set to true.
// path: The path were the yaml to apply is located
// excludes: The list of yamls to exclude
// recursive: If true all yamls in the path directory and sub-directories will be applied
// it excludes the assets named in the excluded array
// it sets the Controller reference if owner and scheme are not nil
//
func (a *Applier) CreateOrUpdateInPath(
	path string,
	excluded []string,
	recursive bool,
	values interface{},
) error {
	us, err := a.templateProcessor.TemplateAssetsInPathUnstructured(
		path,
		excluded,
		recursive,
		values)

	if err != nil {
		return err
	}
	return a.CreateOrUpdates(us)
}

//CreateInPath creates the assets found in the path and
// subpath if recursive is set to true.
// path: The path were the yaml to apply is located
// excludes: The list of yamls to exclude
// recursive: If true all yamls in the path directory and sub-directories will be applied
// it excludes the assets named in the excluded array
// it sets the Controller reference if owner and scheme are not nil
//
func (a *Applier) CreateInPath(
	path string,
	excluded []string,
	recursive bool,
	values interface{},
) error {
	us, err := a.templateProcessor.TemplateAssetsInPathUnstructured(
		path,
		excluded,
		recursive,
		values)

	if err != nil {
		return err
	}
	return a.Creates(us)
}

//UpdateInPath creates or updates the assets found in the path and
// subpath if recursive is set to true.
// path: The path were the yaml to apply is located
// excludes: The list of yamls to exclude
// recursive: If true all yamls in the path directory and sub-directories will be applied
// it excludes the assets named in the excluded array
// it sets the Controller reference if owner and scheme are not nil
//
func (a *Applier) UpdateInPath(
	path string,
	excluded []string,
	recursive bool,
	values interface{},
) error {
	us, err := a.templateProcessor.TemplateAssetsInPathUnstructured(
		path,
		excluded,
		recursive,
		values)

	if err != nil {
		return err
	}
	return a.Updates(us)
}

// Deprecated: Please use another CreateOrUpdateInPath or CreateOrUpdateResouces
// methods with a YamlStringReader
// CreateOrUpdateAssets create or update all resources defined in the assets.
// The asserts are separated by the delimiter (ie: "---" for yamls)
// However the assets is a parameter it requires a reader to define the ToJSON method.
// Please use another method with a YamlStringReader
func (a *Applier) CreateOrUpdateAssets(
	assets []byte,
	values interface{},
	delimiter string,
) error {
	us, err := a.templateProcessor.TemplateBytesUnstructured(assets, values, delimiter)
	if err != nil {
		return err
	}
	return a.CreateOrUpdates(us)
}

//CreateOrUpdateResources creates or update resources
//given an array of resources name
func (a *Applier) CreateOrUpdateResources(
	assetNames []string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateResources(assetNames, values)
	if err != nil {
		return err
	}
	us, err := a.templateProcessor.bytesArrayToUnstructured(b)
	if err != nil {
		return err
	}
	return a.CreateOrUpdates(us)
}

//CreateResources creates resources
//given an array of resources name
func (a *Applier) CreateResources(
	assetNames []string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateResources(assetNames, values)
	if err != nil {
		return err
	}
	us, err := a.templateProcessor.bytesArrayToUnstructured(b)
	if err != nil {
		return err
	}
	return a.Creates(us)
}

//UpdateResources update resources
//given an array of resources name
func (a *Applier) UpdateResources(
	assetNames []string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateResources(assetNames, values)
	if err != nil {
		return err
	}
	us, err := a.templateProcessor.bytesArrayToUnstructured(b)
	if err != nil {
		return err
	}
	return a.Updates(us)
}

//Deprecated: Use CreateOrUpdateResource
//CreateorUpdateAsset create or updates an asset
func (a *Applier) CreateOrUpdateAsset(
	assetName string,
	values interface{},
) error {
	return a.CreateOrUpdateResource(assetName, values)
}

//CreateorUpdateAsset create or updates an asset
func (a *Applier) CreateOrUpdateResource(
	assetName string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateResource(assetName, values)
	if err != nil {
		return err
	}
	u, err := a.templateProcessor.BytesToUnstructured(b)
	if err != nil {
		return err
	}
	return a.CreateOrUpdate(u)
}

//CreateResource create an asset
func (a *Applier) CreateResource(
	assetName string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateResource(assetName, values)
	if err != nil {
		return err
	}
	u, err := a.templateProcessor.BytesToUnstructured(b)
	if err != nil {
		return err
	}
	return a.Create(u)
}

//UpdateResource updates an asset
func (a *Applier) UpdateResource(
	assetName string,
	values interface{},
) error {
	b, err := a.templateProcessor.TemplateAsset(assetName, values)
	if err != nil {
		return err
	}
	u, err := a.templateProcessor.BytesToUnstructured(b)
	if err != nil {
		return err
	}
	return a.Update(u)
}

//CreateOrUpdates an array of unstructured.Unstructured
func (a *Applier) CreateOrUpdates(
	us []*unstructured.Unstructured,
) error {
	//Create the unstructured items if they don't exist yet
	for _, u := range us {
		err := a.CreateOrUpdate(u)
		if err != nil {
			return err
		}
	}
	return nil
}

//Creates create resources from an array of unstructured.Unstructured
func (a *Applier) Creates(
	us []*unstructured.Unstructured,
) error {
	//Create the unstructured items if they don't exist yet
	for _, u := range us {
		err := a.Create(u)
		if err != nil {
			return err
		}
	}
	return nil
}

//Updates updates resources from an array of unstructured.Unstructured
func (a *Applier) Updates(
	us []*unstructured.Unstructured,
) error {
	//Update the unstructured items if they don't exist yet
	for _, u := range us {
		err := a.Update(u)
		if err != nil {
			return err
		}
	}
	return nil
}

//CreateOrUpdate creates or updates an unstructured object.
//It will returns an error if it failed and also if it needs to update the object
//and the applier.Merger is not defined.
func (a *Applier) CreateOrUpdate(
	u *unstructured.Unstructured,
) error {

	klog.V(2).Info("Create or update: ",
		" Kind: ", u.GetKind(),
		" Name: ", u.GetName(),
		" Namespace: ", u.GetNamespace())

	//Check if already exists
	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(u.GroupVersionKind())
	errGet := a.client.Get(context.TODO(), types.NamespacedName{Name: u.GetName(), Namespace: u.GetNamespace()}, current)
	if errGet != nil {
		if errors.IsNotFound(errGet) {
			klog.V(2).Info("Create: ",
				" Kind: ", u.GetKind(),
				" Name: ", u.GetName(),
				" Namespace: ", u.GetNamespace())
			return a.Create(u)
		} else {
			klog.Error(errGet, "Error while create", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
			return errGet
		}
	} else {
		klog.V(2).Info("Update:",
			" Kind: ", current.GetKind(),
			" Name: ", current.GetName(),
			" Namespace: ", current.GetNamespace())
		return a.Update(u)
	}
}

//Create creates an unstructured object.
func (a *Applier) Create(
	u *unstructured.Unstructured,
) error {

	klog.V(2).Info("Create: ",
		" Kind: ", u.GetKind(),
		" Name: ", u.GetName(),
		" Namespace: ", u.GetNamespace())
	//Set controller ref
	err := a.setControllerReference(u)
	if err != nil {
		return err
	}
	var clientCreateOptions []client.CreateOption
	if a.applierOptions != nil {
		clientCreateOptions = a.applierOptions.ClientCreateOption
	}
	createOptions := &client.CreateOptions{}
	clientCreateOption := createOptions.ApplyOptions(clientCreateOptions)
	err = retry.OnError(*a.applierOptions.Backoff, func(err error) bool {
		return err != nil
	}, func() error {
		return a.client.Create(context.TODO(), u, clientCreateOption)
	})
	if err != nil {
		klog.Error(err, "Unable to create: ",
			" Kind: ", u.GetKind(),
			" Name: ", u.GetName(),
			" Namespace: ", u.GetNamespace())
		return err
	}

	return nil
}

//Update updates an unstructured object.
//It will returns an error if it failed and also if it needs to update the object
//and the applier.Merger is not defined.
func (a *Applier) Update(
	u *unstructured.Unstructured,
) error {

	klog.V(2).Info("Create: ",
		" Kind: ", u.GetKind(),
		" Name: ", u.GetName(),
		" Namespace: ", u.GetNamespace())
	//Set controller ref
	err := a.setControllerReference(u)
	if err != nil {
		return err
	}

	//Check if already exists
	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(u.GroupVersionKind())
	errGet := a.client.Get(context.TODO(), types.NamespacedName{Name: u.GetName(), Namespace: u.GetNamespace()}, current)
	if errGet != nil {
		klog.Error(errGet, "Error while update: ",
			" Kind: ", u.GetKind(),
			" Name: ", u.GetName(),
			" Namespace: ", u.GetNamespace())
		return errGet
	} else {
		if a.merger == nil {
			return fmt.Errorf("Unable to update %s/%s of Kind %s the merger is nil",
				current.GetKind(),
				current.GetNamespace(),
				current.GetName())
		}
		future, update := a.merger(current, u)
		if update {
			var clientUpdateOptions []client.UpdateOption
			if a.applierOptions != nil {
				clientUpdateOptions = a.applierOptions.ClientUpdateOption
			}
			updatedOptions := &client.UpdateOptions{}
			clientUpdateOption := updatedOptions.ApplyOptions(clientUpdateOptions)
			err = retry.OnError(*a.applierOptions.Backoff, func(err error) bool {
				return err != nil
			}, func() error {
				return a.client.Update(context.TODO(), future, clientUpdateOption)
			})
			if err != nil {
				klog.Error(err, "Error while update: ",
					" Kind: ", u.GetKind(),
					" Name: ", u.GetName(),
					" Namespace: ", u.GetNamespace())
				return err
			}
		} else {
			klog.V(2).Info("No update needed")
		}
	}
	return nil

}

func (a *Applier) setControllerReference(
	u *unstructured.Unstructured,
) error {
	if a.owner != nil && a.scheme != nil {
		if err := controllerutil.SetControllerReference(a.owner, u, a.scheme); err != nil {
			klog.Error(err, "Failed to SetControllerReference: ",
				" Name: ", u.GetName(),
				" Namespace: ", u.GetNamespace())
			return err
		}
	}
	return nil
}
