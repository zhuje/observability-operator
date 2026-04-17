package uiplugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	osv1 "github.com/openshift/api/console/v1"
	osv1alpha1 "github.com/openshift/api/console/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

type UIPluginInfo struct {
	Image                      string
	Korrel8rImage              string
	HealthAnalyzerImage        string
	LokiServiceNames           map[string]string
	TempoServiceNames          map[string]string
	Name                       string
	ConsoleName                string
	DisplayName                string
	ExtraArgs                  []string
	LegacyProxies              []osv1alpha1.ConsolePluginProxy
	Proxies                    []osv1.ConsolePluginProxy
	Role                       *rbacv1.Role
	RoleBinding                *rbacv1.RoleBinding
	ClusterRoles               []*rbacv1.ClusterRole
	ClusterRoleBindings        []*rbacv1.ClusterRoleBinding
	ConfigMap                  *corev1.ConfigMap
	ResourceNamespace          string
	PersesImage                string
	AreMonitoringFeatsDisabled bool
	TLSMinVersion              string
	TLSCiphers                 []string
}

var pluginTypeToConsoleName = map[uiv1alpha1.UIPluginType]string{
	uiv1alpha1.TypeDashboards:           "console-dashboards-plugin",
	uiv1alpha1.TypeTroubleshootingPanel: "troubleshooting-panel-console-plugin",
	uiv1alpha1.TypeDistributedTracing:   "distributed-tracing-console-plugin",
	uiv1alpha1.TypeLogging:              "logging-view-plugin",
	uiv1alpha1.TypeMonitoring:           "monitoring-console-plugin",
}

// translateOpenSSLToGoNameMap maps OpenSSL cipher names (used by OpenShift TLS profiles)
// to Go TLS package cipher names (expected by UI plugins)
// Based on official cipher suite specifications with hex codes for reference
func translateOpenSSLToGoNameMap() map[string]string {
	return map[string]string{
		// TLS 1.3 Mozilla Modern compatibility
		"TLS_AES_128_GCM_SHA256":          "TLS_AES_128_GCM_SHA256",          // 0x13,0x01
		"TLS_AES_256_GCM_SHA384":          "TLS_AES_256_GCM_SHA384",          // 0x13,0x02
		"TLS_CHACHA20_POLY1305_SHA256":    "TLS_CHACHA20_POLY1305_SHA256",    // 0x13,0x03

		// TLS 1.2 Mozilla Intermediate compatibility  
		"ECDHE-ECDSA-AES128-GCM-SHA256":   "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",   // 0xC0,0x2B (HTTP/2 required)
		"ECDHE-RSA-AES128-GCM-SHA256":     "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",     // 0xC0,0x2F (HTTP/2 required)
		"ECDHE-ECDSA-AES256-GCM-SHA384":   "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",   // 0xC0,0x2C
		"ECDHE-RSA-AES256-GCM-SHA384":     "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",     // 0xC0,0x30
		"ECDHE-ECDSA-CHACHA20-POLY1305":   "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256", // 0xCC,0xA9
		"ECDHE-RSA-CHACHA20-POLY1305":     "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",   // 0xCC,0xA8
	}
}

// translateCiphers converts OpenSSL cipher names to Go TLS package format
// to ensure compatibility between OpenShift TLS profiles and Go-based UI plugins
func translateCiphers(opensslCiphers []string) []string {
	translationMap := translateOpenSSLToGoNameMap()
	var goFormatCiphers []string

	for _, cipher := range opensslCiphers {
		if translatedName, ok := translationMap[cipher]; ok {
			goFormatCiphers = append(goFormatCiphers, translatedName)
		} else {
			// Check if this might be a DHE cipher that should be ECDHE
			if strings.HasPrefix(cipher, "DHE-RSA-") {
				// Potential typo: DHE should be ECDHE for modern profiles
				ecdheVersion := strings.Replace(cipher, "DHE-RSA-", "ECDHE-RSA-", 1)
				if translatedName, ok := translationMap[ecdheVersion]; ok {
					goFormatCiphers = append(goFormatCiphers, translatedName)
					continue
				}
			}
			// Assume it's already in Go format or unknown (let plugin handle)
			goFormatCiphers = append(goFormatCiphers, cipher)
		}
	}

	return goFormatCiphers
}

func PluginInfoBuilder(ctx context.Context, k client.Client, dk dynamic.Interface, plugin *uiv1alpha1.UIPlugin, pluginConf UIPluginsConfiguration, compatibilityInfo CompatibilityEntry, clusterVersion string, logger logr.Logger) (*UIPluginInfo, error) {
	image := pluginConf.Images[compatibilityInfo.ImageKey]
	if image == "" {
		return nil, fmt.Errorf("no image provided for plugin type %s with key %s", plugin.Spec.Type, compatibilityInfo.ImageKey)
	}

	namespace := pluginConf.ResourcesNamespace

	var pluginInfo *UIPluginInfo
	var err error

	switch plugin.Spec.Type {
	case uiv1alpha1.TypeDashboards:
		pluginInfo, err = createDashboardsPluginInfo(plugin, namespace, plugin.Name, image)
		if err != nil {
			return nil, err
		}

	case uiv1alpha1.TypeTroubleshootingPanel:
		pluginInfo, err = createTroubleshootingPanelPluginInfo(plugin, namespace, plugin.Name, image, []string{})
		if err != nil {
			return nil, err
		}

		pluginInfo.Korrel8rImage = pluginConf.Images["korrel8r"]
		pluginInfo.LokiServiceNames[OpenshiftLoggingNs], err = getLokiServiceName(ctx, k, OpenshiftLoggingNs)
		if err != nil {
			return nil, err
		}

		pluginInfo.LokiServiceNames[OpenshiftNetobservNs], err = getLokiServiceName(ctx, k, OpenshiftNetobservNs)
		if err != nil {
			return nil, err
		}

		pluginInfo.TempoServiceNames[OpenshiftTracingNs], err = getTempoServiceName(ctx, k, OpenshiftTracingNs)
		if err != nil {
			return nil, err
		}

	case uiv1alpha1.TypeDistributedTracing:
		pluginInfo, err = createDistributedTracingPluginInfo(plugin, namespace, plugin.Name, image, []string{})
		if err != nil {
			return nil, err
		}

	case uiv1alpha1.TypeLogging:
		pluginInfo, err = createLoggingPluginInfo(plugin, namespace, plugin.Name, image, compatibilityInfo.Features, ctx, dk, logger, pluginConf.Images["korrel8r"])
		if err != nil {
			return nil, err
		}

	case uiv1alpha1.TypeMonitoring:
		pluginInfo, err = createMonitoringPluginInfo(plugin, namespace, plugin.Name, image, compatibilityInfo.Features, clusterVersion, pluginConf.Images["health-analyzer"], pluginConf.Images["perses"])
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("plugin type not supported: %s", plugin.Spec.Type)
	}

	if compatibilityInfo.SupportsTLSProfile {
		pluginInfo.TLSMinVersion = string(pluginConf.TLSProfile.MinTLSVersion)
		// Translate OpenSSL cipher names to Go format before passing to plugin
		pluginInfo.TLSCiphers = translateCiphers(pluginConf.TLSProfile.Ciphers)
	} else {
		logger.Info("TLS profile not applied: plugin image does not support TLS profile flags",
			"plugin", plugin.Name,
			"pluginType", plugin.Spec.Type,
			"imageKey", compatibilityInfo.ImageKey)
	}

	return pluginInfo, err
}
