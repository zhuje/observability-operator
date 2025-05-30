package operator

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	stackctrl "github.com/rhobs/observability-operator/pkg/controllers/monitoring/monitoring-stack"
	tqctrl "github.com/rhobs/observability-operator/pkg/controllers/monitoring/thanos-querier"
	opctrl "github.com/rhobs/observability-operator/pkg/controllers/operator"
	uictrl "github.com/rhobs/observability-operator/pkg/controllers/uiplugin"
)

const (
	// NOTE: The instance selector label is hardcoded in static assets.
	// Any change to that must be reflected here as well
	instanceSelector = "app.kubernetes.io/managed-by=observability-operator"

	ObservabilityOperatorName = "observability-operator"

	// The mount path for the serving certificate seret is hardcoded in the
	// static assets.
	tlsMountPath = "/etc/tls/private"
)

// Operator embeds a manager and a serving certificate controller (for
// OpenShift installations).
type Operator struct {
	manager               manager.Manager
	servingCertController *dynamiccertificates.DynamicServingCertificateController
	clientCAController    *dynamiccertificates.ConfigMapCAController
}

type OpenShiftFeatureGates struct {
	Enabled bool `json:"enabled,omitempty"`
}

type FeatureGates struct {
	OpenShift OpenShiftFeatureGates `json:"openshift,omitempty"`
}

type OperatorConfiguration struct {
	Namespace       string
	MetricsAddr     string
	HealthProbeAddr string
	Prometheus      stackctrl.PrometheusConfiguration
	Alertmanager    stackctrl.AlertmanagerConfiguration
	ThanosSidecar   stackctrl.ThanosConfiguration
	ThanosQuerier   tqctrl.ThanosConfiguration
	UIPlugins       uictrl.UIPluginsConfiguration
	FeatureGates    FeatureGates
}

func WithNamespace(ns string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.Namespace = ns
		oc.UIPlugins.ResourcesNamespace = ns
	}
}

func WithPrometheusImage(image string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.Prometheus.Image = image
	}
}

func WithAlertmanagerImage(image string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.Alertmanager.Image = image
	}
}

func WithThanosSidecarImage(image string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.ThanosSidecar.Image = image
	}
}

func WithThanosQuerierImage(image string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.ThanosQuerier.Image = image
	}
}

func WithMetricsAddr(addr string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.MetricsAddr = addr
	}
}

func WithHealthProbeAddr(addr string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.HealthProbeAddr = addr
	}
}

func WithUIPluginImages(images map[string]string) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.UIPlugins.Images = images
	}
}

func WithFeatureGates(featureGates FeatureGates) func(*OperatorConfiguration) {
	return func(oc *OperatorConfiguration) {
		oc.FeatureGates = featureGates
	}
}

func NewOperatorConfiguration(opts ...func(*OperatorConfiguration)) *OperatorConfiguration {
	cfg := &OperatorConfiguration{}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func New(ctx context.Context, cfg *OperatorConfiguration) (*Operator, error) {
	restConfig := ctrl.GetConfigOrDie()
	scheme := NewScheme(cfg)

	metricsOpts := metricsserver.Options{
		BindAddress: cfg.MetricsAddr,
	}

	var (
		clientCAController    *dynamiccertificates.ConfigMapCAController
		servingCertController *dynamiccertificates.DynamicServingCertificateController
	)
	if cfg.FeatureGates.OpenShift.Enabled {
		// When running in OpenShift, the server uses HTTPS thanks to the
		// service CA operator.
		certFile := filepath.Join(tlsMountPath, "tls.crt")
		keyFile := filepath.Join(tlsMountPath, "tls.key")

		// Wait for the files to be mounted into the container.
		var pollErr error
		err := wait.PollUntilContextTimeout(ctx, time.Second, 30*time.Second, true, func(ctx context.Context) (bool, error) {
			for _, f := range []string{certFile, keyFile} {
				if _, err := os.Stat(f); err != nil {
					pollErr = err
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", err, pollErr)
		}

		// DynamicCertKeyPairContent automatically reloads the certificate and key from disk.
		certKeyProvider, err := dynamiccertificates.NewDynamicServingContentFromFiles("serving-cert", certFile, keyFile)
		if err != nil {
			return nil, err
		}
		if err := certKeyProvider.RunOnce(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize cert/key content: %w", err)
		}

		kubeClient, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			return nil, err
		}

		clientCAController, err = dynamiccertificates.NewDynamicCAFromConfigMapController(
			"client-ca",
			metav1.NamespaceSystem,
			"extension-apiserver-authentication",
			"client-ca-file",
			kubeClient,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize client CA controller: %w", err)
		}

		// Only log the events emitted by the certificate controller for now
		// because the controller generates invalid events rejected by the
		// Kubernetes API when used with DynamicServingContentFromFiles.
		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartLogging(func(format string, args ...interface{}) {
			ctrl.Log.WithName("events").Info(fmt.Sprintf(format, args...))
		})

		servingCertController = dynamiccertificates.NewDynamicServingCertificateController(
			&tls.Config{
				ClientAuth: tls.RequireAndVerifyClientCert,
			},
			clientCAController,
			certKeyProvider,
			nil,
			record.NewEventRecorderAdapter(
				eventBroadcaster.NewRecorder(scheme, v1.EventSource{Component: "observability-operator"}),
			),
		)
		if err := servingCertController.RunOnce(); err != nil {
			return nil, fmt.Errorf("failed to initialize serving certificate controller: %w", err)
		}

		clientCAController.AddListener(servingCertController)
		certKeyProvider.AddListener(servingCertController)

		metricsOpts.SecureServing = true
		metricsOpts.TLSOpts = []func(*tls.Config){
			func(c *tls.Config) {
				c.GetConfigForClient = servingCertController.GetConfigForClient
			},
		}
	}

	mgr, err := ctrl.NewManager(
		restConfig,
		ctrl.Options{
			Scheme:                 scheme,
			Metrics:                metricsOpts,
			HealthProbeBindAddress: cfg.HealthProbeAddr,
			PprofBindAddress:       "127.0.0.1:8083",
		})
	if err != nil {
		return nil, fmt.Errorf("unable to create manager: %w", err)
	}

	if err := stackctrl.RegisterWithManager(mgr, stackctrl.Options{
		InstanceSelector: instanceSelector,
		Prometheus:       cfg.Prometheus,
		Alertmanager:     cfg.Alertmanager,
		Thanos:           cfg.ThanosSidecar,
	}); err != nil {
		return nil, fmt.Errorf("unable to register monitoring stack controller: %w", err)
	}

	if err := tqctrl.RegisterWithManager(mgr, tqctrl.Options{Thanos: cfg.ThanosQuerier}); err != nil {
		return nil, fmt.Errorf("unable to register the thanos querier controller with the manager: %w", err)
	}

	if cfg.FeatureGates.OpenShift.Enabled {
		if err := uictrl.RegisterWithManager(mgr, uictrl.Options{PluginsConf: cfg.UIPlugins}); err != nil {
			return nil, fmt.Errorf("unable to register observability-ui-plugin controller: %w", err)
		}
	} else {
		setupLog := ctrl.Log.WithName("setup")
		setupLog.Info("OpenShift feature gate is disabled, UIPlugins are not enabled")
	}

	if cfg.FeatureGates.OpenShift.Enabled {
		if err := opctrl.RegisterWithManager(mgr, cfg.Namespace); err != nil {
			return nil, fmt.Errorf("unable to register operator controller: %w", err)
		}
	} else {
		setupLog := ctrl.Log.WithName("setup")
		setupLog.Info("OpenShift feature gate is disabled, Operator controller is not enabled")
	}

	if err := mgr.AddHealthzCheck("health probe", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to add health probe: %w", err)
	}

	return &Operator{
		manager:               mgr,
		servingCertController: servingCertController,
		clientCAController:    clientCAController,
	}, nil
}

func (o *Operator) Start(ctx context.Context) error {
	if o.clientCAController != nil {
		go o.clientCAController.Run(ctx, 1)
	}

	if o.servingCertController != nil {
		go o.servingCertController.Run(1, ctx.Done())
	}

	if err := o.manager.Start(ctx); err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	return nil
}

func (o *Operator) GetClient() client.Client {
	return o.manager.GetClient()
}
