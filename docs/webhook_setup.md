# Webhook

This package contains `WebHookWireUp` Object, which provides a convenient way to
wire up a `ValidatingWebhookConfiguration` webhook to the `manager` of your
operator.

In addition, the `Webhook` package also provide helpers(`certificate.go`) to help your, 
1. create a key pair at `CertDir`
2. create CA cert, which will be injected to the `ValidatingWebhookConfiguration`

# How to use it

1. Create a new `WebHookWireUp` instance, via `NewWebHookWireUp()` 
2. Call `WebHookWireUp.Attach()` to link the `manager` and `webhook`. At this
   point the `webhook` had the validation Handler injected.
3. Use a goroutine to run the
   ```
   WebHookWireUp.WireUpWebhookSupplymentryResource(caCert []byte, gvk
   schema.GroupVersionKind, ops []admissionv1.OperationType)
   ```
   Which will check if there are:
	a. webhookName + "-svc"
	b. webhookName + '.validator'
	
	If either of the above exists, we will create it and set the owner to the pod
	who's running the manager. We will use the `podNamespaceEnvVar` and 
	`deployLabelEnvVar` to find the deployment.
	
	The `ValidatingWebhookConfiguration` will run the handle logic on the
	resource specified by `gvk` when the operations/events, defined at `ops`,
	happens.
	
As the end result, you should find a `service` running in the pod's namespace,
also a `ValidatingWebhookConfiguration` point to the `service`. When users
operate on the watching resources of given operations, the validation logic will
be triggered. Here's an [exmple](https://github.com/open-cluster-management/multicloud-operators-channel/blob/master/cmd/manager/exec/manager.go#L217)

# How to customize 
As you can see for the `NewWebHookWireUp()`, it expects the optional function.
In this way, you can override all the exported fields. For example, to change
the `CertDir` field, you can do the following. 

```go
	wbhCertDir := func(w *chWebhook.WebHookWireUp) {
		w.CertDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
	}

	wiredWebhook, err := chWebhook.NewWebHookWireUp(mgr, sig, wbhCertDir)
	if err != nil {
		logger.Error(err, "failed to initial wire up webhook")
		os.Exit(exitCode)
	}
```

You can also refer to the `wireupwebhook_test.go` for some sample code. 
