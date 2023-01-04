// Copyright Contributors to the Open Cluster Management project

package webhook

import (
	"context"
	"fmt"
	"os"
	"strings"

	gerr "github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

	admissionv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	tlsCrt = "tls.crt"
	tlsKey = "tls.key"
)

type WireUp struct {
	mgr  manager.Manager
	stop <-chan struct{}

	Server  *webhook.Server
	Handler webhook.AdmissionHandler
	CertDir string
	logr.Logger

	WebhookName string
	WebHookPort int

	WebHookeSvcKey     types.NamespacedName
	WebHookServicePort int

	ValidtorPath string

	DeployLabel        string
	DeploymentSelector map[string]string
}

func NewWireUp(mgr manager.Manager, stop <-chan struct{},
	opts ...func(*WireUp)) (*WireUp, error) {

	WebhookName := "channels-apps-open-cluster-management-webhook"

	deployLabelEnvVar := "DEPLOYMENT_LABEL"
	deploymentLabel, err := findEnvVariable(deployLabelEnvVar)

	if err != nil {
		return nil, gerr.Wrap(err, "failed to create a new webhook wireup")
	}

	deploymentSelectKey := "app"

	deploymentSelector := map[string]string{deploymentSelectKey: deploymentLabel}

	podNamespaceEnvVar := "POD_NAMESPACE"

	podNamespace, err := findEnvVariable(podNamespaceEnvVar)
	if err != nil {
		return nil, gerr.Wrap(err, "failed to create a new webhook wireup")
	}

	wireUp := &WireUp{
		mgr:    mgr,
		stop:   stop,
		Server: mgr.GetWebhookServer(),
		Logger: logf.NullLogger{},

		WebhookName:        WebhookName,
		WebHookPort:        9443,
		WebHookServicePort: 443,

		ValidtorPath:       "/duplicate-validation",
		WebHookeSvcKey:     types.NamespacedName{Name: GetWebHookServiceName(WebhookName), Namespace: podNamespace},
		DeployLabel:        deploymentLabel,
		DeploymentSelector: deploymentSelector,
	}

	for _, op := range opts {
		op(wireUp)
	}

	return wireUp, nil
}

func GetValidatorName(wbhName string) string {
	//domain style, separate by dots
	return strings.ReplaceAll(wbhName, "-", ".") + ".validator"
}

func GetWebHookServiceName(wbhName string) string {
	//k8s resrouce naming, separate by '-'
	return wbhName + "-svc"
}

func (w *WireUp) Attach() ([]byte, error) {
	w.Server.Port = w.WebHookPort

	w.Logger.Info("registering webhooks to the webhook server")
	w.Server.Register(w.ValidtorPath, &webhook.Admission{Handler: w.Handler})

	return GenerateWebhookCerts(w.CertDir, w.WebHookeSvcKey.Namespace, w.WebHookeSvcKey.Name)
}

//assuming we have a service set up for the webhook, and the service is linking
//to a secret which has the CA
func (w *WireUp) WireUpWebhookSupplymentryResource(caCert []byte, gvk schema.GroupVersionKind, ops []admissionv1.OperationType) {
	w.Logger.Info("entry wire up webhook")
	defer w.Logger.Info("exit wire up webhook ")

	if !w.mgr.GetCache().WaitForCacheSync(w.stop) {
		w.Logger.Error(gerr.New("cache not started"), "failed to start up cache")
	}

	w.Logger.Info("cache is ready to consume")

	if err := w.createWebhookService(); err != nil {
		w.Logger.Error(err, "failed to wire up webhook with kube")
		os.Exit(1)
	}

	if err := w.createOrUpdateValiatingWebhookConfiguration(caCert, gvk, ops); err != nil {
		w.Logger.Error(err, "failed to wire up webhook with kube")
		os.Exit(1)
	}
}

func findEnvVariable(envName string) (string, error) {
	val, found := os.LookupEnv(envName)
	if !found {
		return "", fmt.Errorf("%s env var is not set", envName)
	}

	return val, nil
}

func (w *WireUp) createWebhookService() error {
	service := &corev1.Service{}

	if err := w.mgr.GetClient().Get(context.TODO(), w.WebHookeSvcKey, service); err != nil {
		if errors.IsNotFound(err) {
			service := newWebhookServiceTemplate(w.WebHookeSvcKey, w.WebHookPort, w.WebHookServicePort, w.DeploymentSelector)

			setOwnerReferences(w.mgr.GetClient(), w.Logger, w.WebHookeSvcKey.Namespace, w.DeployLabel, service)

			if err := w.mgr.GetClient().Create(context.TODO(), service); err != nil {
				return err
			}

			w.Logger.Info(fmt.Sprintf("Create %s service", w.WebHookeSvcKey.String()))

			return nil
		}
	}

	w.Logger.Info(fmt.Sprintf("%s service is found", w.WebHookeSvcKey.String()))

	return nil
}

func (w *WireUp) createOrUpdateValiatingWebhookConfiguration(ca []byte, gvk schema.GroupVersionKind,
	ops []admissionv1.OperationType) error {
	validator := &admissionv1.ValidatingWebhookConfiguration{}
	key := types.NamespacedName{Name: GetValidatorName(w.WebhookName)}

	validatorName := GetValidatorName(w.WebhookName)
	if err := w.mgr.GetClient().Get(context.TODO(), key, validator); err != nil {
		if errors.IsNotFound(err) {
			cfg := newValidatingWebhookCfg(w.WebHookeSvcKey, validatorName, w.ValidtorPath, ca, gvk, ops)

			setOwnerReferences(w.mgr.GetClient(), w.Logger, w.WebHookeSvcKey.Namespace, w.DeployLabel, cfg)

			if err := w.mgr.GetClient().Create(context.TODO(), cfg); err != nil {
				return gerr.Wrap(err, fmt.Sprintf("Failed to create validating webhook %s", validatorName))
			}

			w.Logger.Info(fmt.Sprintf("Create validating webhook %s", validatorName))

			return nil
		}
	}

		// since we only have 1 webhook over here, override the original one is ok
	validator.Webhooks = []admissionv1.ValidatingWebhook{
		{
			ClientConfig: admissionv1.WebhookClientConfig{
				Service:  &admissionv1.ServiceReference{Namespace: w.WebHookeSvcKey.Namespace},
				CABundle: ca,
			},
		},
	}

	if err := w.mgr.GetClient().Update(context.TODO(), validator); err != nil {
		return gerr.Wrap(err, fmt.Sprintf("Failed to update validating webhook %s", validatorName))
	}

	w.Logger.Info(fmt.Sprintf("Update validating webhook %s", validatorName))

	return nil
}

func setOwnerReferences(c client.Client, logger logr.Logger, deployNs string, deployLabel string, obj metav1.Object) {
	key := types.NamespacedName{Name: deployLabel, Namespace: deployNs}
	owner := &appsv1.Deployment{}

	if err := c.Get(context.TODO(), key, owner); err != nil {
		logger.Error(err, fmt.Sprintf("Failed to set owner references for %s", obj.GetName()))
		return
	}

	obj.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(owner, owner.GetObjectKind().GroupVersionKind())})
}

func newWebhookServiceTemplate(svcKey types.NamespacedName, webHookPort, webHookServicePort int, deploymentSelector map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcKey.Name,
			Namespace: svcKey.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       int32(webHookServicePort),
					TargetPort: intstr.FromInt(webHookPort),
				},
			},
			Selector: deploymentSelector,
		},
	}
}

func newValidatingWebhookCfg(svcKey types.NamespacedName, validatorName, path string, ca []byte,
	gvk schema.GroupVersionKind, ops []admissionv1.OperationType) *admissionv1.ValidatingWebhookConfiguration {
	fail := admissionv1.Fail
	side := admissionv1.SideEffectClassNone

	return &admissionv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: validatorName,
		},

		Webhooks: []admissionv1.ValidatingWebhook{{
			Name:                    validatorName,
			AdmissionReviewVersions: []string{"v1beta1"},
			SideEffects:             &side,
			FailurePolicy:           &fail,
			ClientConfig: admissionv1.WebhookClientConfig{
				Service: &admissionv1.ServiceReference{
					Name:      svcKey.Name,
					Namespace: svcKey.Namespace,
					Path:      &path,
				},
				CABundle: ca,
			},
			Rules: []admissionv1.RuleWithOperations{{
				Rule: admissionv1.Rule{
					APIGroups:   []string{gvk.Group},
					APIVersions: []string{gvk.Version},
					Resources:   []string{gvk.Kind},
				},
				Operations: ops,
			},
			}},
		},
	}
}
