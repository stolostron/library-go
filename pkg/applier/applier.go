package applier

import (
	"context"
	goerr "errors"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
) (*Applier, error) {
	if templateProcessor == nil {
		return nil, goerr.New("applier is nil")
	}
	if client == nil {
		return nil, goerr.New("client is nil")
	}
	return &Applier{
		templateProcessor: templateProcessor,
		client:            client,
		owner:             owner,
		scheme:            scheme,
		merger:            merger,
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
	if spec, ok := new.Object["spec"]; ok &&
		!reflect.DeepEqual(spec, current.Object["spec"]) {
		update = true
		current.Object["spec"] = spec
	}
	if rules, ok := new.Object["rules"]; ok &&
		!reflect.DeepEqual(rules, current.Object["rules"]) {
		update = true
		current.Object["rules"] = rules
	}
	if roleRef, ok := new.Object["roleRef"]; ok &&
		!reflect.DeepEqual(roleRef, current.Object["roleRef"]) {
		update = true
		current.Object["roleRef"] = roleRef
	}
	if subjects, ok := new.Object["subjects"]; ok &&
		!reflect.DeepEqual(subjects, current.Object["subjects"]) {
		update = true
		current.Object["subjects"] = subjects
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
	return a.CreateArrayOfUnstructuted(us)
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
	return a.UpdateArrayOfUnstructured(us)
}

//CreateOrUpdateAssets create or update all resources defined in the assets.
//The asserts are separated by the delimiter (ie: "---" for yamls)
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

//CreatAssets create all resources defined in the assets.
//The asserts are separated by the delimiter (ie: "---" for yamls)
func (a *Applier) CreateAssets(
	assets []byte,
	values interface{},
	delimiter string,
) error {
	us, err := a.templateProcessor.TemplateBytesUnstructured(assets, values, delimiter)
	if err != nil {
		return err
	}
	return a.CreateArrayOfUnstructuted(us)
}

//UpdateAssets update all resources defined in the assets.
//The asserts are separated by the delimiter (ie: "---" for yamls)
func (a *Applier) UpdateAssets(
	assets []byte,
	values interface{},
	delimiter string,
) error {
	us, err := a.templateProcessor.TemplateBytesUnstructured(assets, values, delimiter)
	if err != nil {
		return err
	}
	return a.UpdateArrayOfUnstructured(us)
}

//CreateorUpdateAsset create or updates an asset
func (a *Applier) CreateOrUpdateAsset(
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
	return a.CreateOrUpdate(u)
}

//CreateAsset create an asset
func (a *Applier) CreateAsset(
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
	return a.Create(u)
}

//UpdateAsset updates an asset
func (a *Applier) UpdateAsset(
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

//CreateArrayOfUnstructured create resources from an array of unstructured.Unstructured
func (a *Applier) CreateArrayOfUnstructuted(
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

//UpdateArrayOfUnstructured updates resources from an array of unstructured.Unstructured
func (a *Applier) UpdateArrayOfUnstructured(
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

	log.V(2).Info("Create or update", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())

	//Check if already exists
	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(u.GroupVersionKind())
	errGet := a.client.Get(context.TODO(), types.NamespacedName{Name: u.GetName(), Namespace: u.GetNamespace()}, current)
	if errGet != nil {
		if errors.IsNotFound(errGet) {
			log.V(2).Info("Create",
				"Kind", current.GetKind(),
				"Name", current.GetName(),
				"Namespace", current.GetNamespace())
			return a.Create(u)
		} else {
			log.Error(errGet, "Error while create", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
			return errGet
		}
	} else {
		log.V(2).Info("Update",
			"Kind", current.GetKind(),
			"Name", current.GetName(),
			"Namespace", current.GetNamespace())
		return a.Update(u)
	}
}

//Create creates an unstructured object.
func (a *Applier) Create(
	u *unstructured.Unstructured,
) error {

	log.V(2).Info("Create ", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
	//Set controller ref
	err := a.setControllerReference(u)
	if err != nil {
		return err
	}

	err = a.client.Create(context.TODO(), u)
	if err != nil {
		log.Error(err, "Unable to create", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
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

	log.V(2).Info("Create ", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
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
		log.Error(errGet, "Error while update", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
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
			err := a.client.Update(context.TODO(), future)
			if err != nil {
				log.Error(err, "Error while update", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
				return err
			}
		} else {
			log.V(2).Info("No update needed")
		}
	}
	return nil

}

func (a *Applier) setControllerReference(
	u *unstructured.Unstructured,
) error {
	if a.owner != nil && a.scheme != nil {
		if err := controllerutil.SetControllerReference(a.owner, u, a.scheme); err != nil {
			log.Error(err, "Failed to SetControllerReference",
				"Name", u.GetName(),
				"Namespace", u.GetNamespace())
			return err
		}
	}
	return nil
}
