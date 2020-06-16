// Copyright (c) 2020 Red Hat, Inc.

package applier

import (
	"context"
	goerr "errors"
	"fmt"

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

//CreateOrUpdates an array of unstructured.Unstructured
func (a *Applier) CreateOrUpdates(us []*unstructured.Unstructured) error {
	//Create the unstructured items if they don't exist yet
	for _, u := range us {
		err := a.CreateOrUpdate(u)
		if err != nil {
			return err
		}
	}
	return nil
}

//CreateOrUpdate creates or updates an unstructured (if they don't exist yet) found in the path and
func (a *Applier) CreateOrUpdate(
	u *unstructured.Unstructured,
) error {

	log.Info("Create or update", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
	//Set controller ref
	if a.owner != nil && a.scheme != nil {
		if err := controllerutil.SetControllerReference(a.owner, u, a.scheme); err != nil {
			log.Error(err, "Failed to SetControllerReference",
				"Name", u.GetName(),
				"Namespace", u.GetNamespace())
			return err
		}
	}

	//Check if already exists
	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(u.GroupVersionKind())
	var errGet error
	errGet = a.client.Get(context.TODO(), types.NamespacedName{Name: u.GetName(), Namespace: u.GetNamespace()}, current)
	if errGet != nil {
		if errors.IsNotFound(errGet) {
			log.Info("Create",
				"Kind", current.GetKind(),
				"Name", current.GetName(),
				"Namespace", current.GetNamespace())
			err := a.client.Create(context.TODO(), u)
			if err != nil {
				log.Error(err, "Unable to create", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
				return err
			}
		} else {
			log.Error(errGet, "Error while create", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
			return errGet
		}
	} else {
		log.Info("Update",
			"Kind", current.GetKind(),
			"Name", current.GetName(),
			"Namespace", current.GetNamespace())
		if a.merger == nil {
			return fmt.Errorf("Unable to update as the merger is nil")
		}
		future, update := a.merger(current, u)
		if update {
			err := a.client.Update(context.TODO(), future)
			if err != nil {
				log.Error(err, "Error while update", "Kind", u.GetKind(), "Name", u.GetName(), "Namespace", u.GetNamespace())
				return err
			}
		} else {
			log.Info("No update needed")
		}
	}
	return nil
}
