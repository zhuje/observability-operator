package uiplugin

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-logr/logr"
	osv1 "github.com/openshift/api/console/v1"
	osv1alpha1 "github.com/openshift/api/console/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

func createMonitoringPluginInfo(plugin *uiv1alpha1.UIPlugin, namespace, name, image string, features []string, logger logr.Logger) (*UIPluginInfo, error) {
	config := plugin.Spec.Monitoring
	if config == nil {
		return nil, fmt.Errorf("monitoring configuration can not be empty for plugin type %s", plugin.Spec.Type)
	}

	// Default service name and namespace for Perses
	persesName := "perses-api-http"
	persesNamespace := "perses-operator"

	persesDashboardsFeatureEnabled := slices.Contains(features, "perses-dashboards")

	logger.Info("6. HelloWorld ", "config.Perses.Name", config.Perses.Name, "config.Perses.Namespace", config.Perses.Namespace)

	// validate UIPlugin CR monitoring properties
	validACMConfig := (config.Alertmanager.Url != "" && config.ThanosQuerier.Url != "")
	validConfig := validACMConfig || persesDashboardsFeatureEnabled

	if !validConfig {
		if config.Alertmanager.Url == "" {
			return nil, fmt.Errorf("alertmanager location can not be empty for plugin type %s", plugin.Spec.Type)
		}
		if config.ThanosQuerier.Url == "" {
			return nil, fmt.Errorf("ThanosQuerier location can not be empty for plugin type %s", plugin.Spec.Type)
		}
	}

	// allow UIPlugin CR to override default perses name and namespace
	if persesDashboardsFeatureEnabled && config.Perses.Name != "" && config.Perses.Namespace != "" {
		persesName = config.Perses.Name
		persesNamespace = config.Perses.Namespace
	}

	logger.Info("6. HelloWorld ", "config.Perses.Name", config.Perses.Name, "config.Perses.Namespace", config.Perses.Namespace)

	pluginInfo := &UIPluginInfo{
		Image:       image,
		Name:        name,
		ConsoleName: "monitoring-console-plugin",
		DisplayName: "Monitoring Console Plugin",
		ExtraArgs: []string{
			fmt.Sprintf("-features=%s", strings.Join(features, ",")),
			fmt.Sprintf("-alertmanager=%s", config.Alertmanager.Url),
			fmt.Sprintf("-thanos-querier=%s", config.ThanosQuerier.Url),
			"-config-path=/opt/app-root/config",
			"-static-path=/opt/app-root/web/dist",
		},
		ResourceNamespace: namespace,
		Proxies: []osv1.ConsolePluginProxy{
			{
				Alias:         "backend",
				Authorization: "UserToken",
				Endpoint: osv1.ConsolePluginProxyEndpoint{
					Type: osv1.ProxyTypeService,
					Service: &osv1.ConsolePluginProxyServiceConfig{
						Name:      name,
						Namespace: namespace,
						Port:      port,
					},
				},
			},
			{
				Alias:         "alertmanager-proxy",
				Authorization: "UserToken",
				Endpoint: osv1.ConsolePluginProxyEndpoint{
					Type: osv1.ProxyTypeService,
					Service: &osv1.ConsolePluginProxyServiceConfig{
						Name:      name,
						Namespace: namespace,
						Port:      9444,
					},
				},
			},
			{
				Alias:         "thanos-proxy",
				Authorization: "UserToken",
				Endpoint: osv1.ConsolePluginProxyEndpoint{
					Type: osv1.ProxyTypeService,
					Service: &osv1.ConsolePluginProxyServiceConfig{
						Name:      name,
						Namespace: namespace,
						Port:      9445,
					},
				},
			},
		},
		LegacyProxies: []osv1alpha1.ConsolePluginProxy{
			{
				Type:      "Service",
				Alias:     "backend",
				Authorize: true,
				Service: osv1alpha1.ConsolePluginProxyServiceConfig{
					Name:      name,
					Namespace: namespace,
					Port:      9443,
				},
			},
			{
				Type:      "Service",
				Alias:     "alertmanager-proxy",
				Authorize: true,
				Service: osv1alpha1.ConsolePluginProxyServiceConfig{
					Name:      name,
					Namespace: namespace,
					Port:      9444,
				},
			},
			{
				Type:      "Service",
				Alias:     "thanos-proxy",
				Authorize: true,
				Service: osv1alpha1.ConsolePluginProxyServiceConfig{
					Name:      name,
					Namespace: namespace,
					Port:      9445,
				},
			},
		},
	}

	if persesDashboardsFeatureEnabled {
		pluginInfo.Proxies = append(pluginInfo.Proxies, osv1.ConsolePluginProxy{
			Alias:         "perses",
			Authorization: "UserToken",
			Endpoint: osv1.ConsolePluginProxyEndpoint{
				Type: osv1.ProxyTypeService,
				Service: &osv1.ConsolePluginProxyServiceConfig{
					Name:      persesName,
					Namespace: persesNamespace,
					Port:      8080,
				},
			},
		})
		pluginInfo.LegacyProxies = append(pluginInfo.LegacyProxies, osv1alpha1.ConsolePluginProxy{
			Type:      "Service",
			Alias:     "perses",
			Authorize: true,
			Service: osv1alpha1.ConsolePluginProxyServiceConfig{
				Name:      persesName,
				Namespace: persesNamespace,
				Port:      8080,
			},
		})
	}

	return pluginInfo, nil
}

func newMonitoringService(name string, namespace string) *corev1.Service {
	annotations := map[string]string{
		"service.beta.openshift.io/serving-cert-secret-name": name,
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      componentLabels(name),
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       9443,
					Name:       "backend",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt32(9443),
				},
				{
					Port:       9444,
					Name:       "alertmanager-proxy",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt32(9444),
				},
				{
					Port:       9445,
					Name:       "thanos-proxy",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt32(9445),
				},
			},
			Selector: componentLabels(name),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}
