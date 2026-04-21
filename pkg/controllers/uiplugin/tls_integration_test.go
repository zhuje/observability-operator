package uiplugin

import (
	"context"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"gotest.tools/v3/assert"
	"sigs.k8s.io/controller-runtime/pkg/log"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

func TestTLSProfileEndToEndFlow(t *testing.T) {
	testCases := []struct {
		name                string
		tlsProfile          configv1.TLSProfileSpec
		pluginType          uiv1alpha1.UIPluginType
		clusterVersion      string
		supportsTLSProfile  bool
		expectedTLSArgs     bool
		expectedMinVersion  string
		expectedCiphers     []string
	}{
		{
			name: "Monitoring plugin v4.19+ with TLS profile",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS12,
				Ciphers:       []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"},
			},
			pluginType:         uiv1alpha1.TypeMonitoring,
			clusterVersion:     "v4.19",
			supportsTLSProfile: true,
			expectedTLSArgs:    true,
			expectedMinVersion: string(configv1.VersionTLS12),
			expectedCiphers:    []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"},
		},
		{
			name: "Monitoring plugin v4.18 (no TLS support)",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS12,
				Ciphers:       []string{"TLS_AES_128_GCM_SHA256"},
			},
			pluginType:         uiv1alpha1.TypeMonitoring,
			clusterVersion:     "v4.18",
			supportsTLSProfile: false,
			expectedTLSArgs:    false,
			expectedMinVersion: "",
			expectedCiphers:    nil,
		},
		{
			name: "Dashboards plugin (no TLS support yet)",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS13,
				Ciphers:       []string{"TLS_AES_256_GCM_SHA384"},
			},
			pluginType:         uiv1alpha1.TypeDashboards,
			clusterVersion:     "v4.19",
			supportsTLSProfile: false,
			expectedTLSArgs:    false,
			expectedMinVersion: "",
			expectedCiphers:    nil,
		},
		{
			name: "Empty TLS profile with supporting plugin",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: "",
				Ciphers:       nil,
			},
			pluginType:         uiv1alpha1.TypeMonitoring,
			clusterVersion:     "v4.19",
			supportsTLSProfile: true,
			expectedTLSArgs:    false, // No args when profile is empty
			expectedMinVersion: "",
			expectedCiphers:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			logger := log.FromContext(ctx)

			// Create mock UIPlugin
			plugin := &uiv1alpha1.UIPlugin{
				Spec: uiv1alpha1.UIPluginSpec{
					Type: tc.pluginType,
				},
			}

			// Add monitoring configuration for monitoring plugins
			if tc.pluginType == uiv1alpha1.TypeMonitoring {
				plugin.Spec.Monitoring = &uiv1alpha1.MonitoringConfig{
					ACM: &uiv1alpha1.AdvancedClusterManagementReference{
						Enabled: true,
						Alertmanager: uiv1alpha1.AlertmanagerReference{
							Url: "https://alertmanager.test.svc:9095",
						},
						ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
							Url: "https://thanos.test.svc:9091",
						},
					},
				}
			}

			// Create UIPluginsConfiguration with TLS profile
			pluginConf := UIPluginsConfiguration{
				TLSProfile: tc.tlsProfile,
				Images: map[string]string{
					"ui-monitoring":     "quay.io/test/monitoring:latest",
					"ui-monitoring-pf5": "quay.io/test/monitoring-pf5:latest",
					"ui-dashboards":     "quay.io/test/dashboards:latest",
					"health-analyzer":   "quay.io/test/health:latest",
					"perses":            "quay.io/test/perses:latest",
				},
				ResourcesNamespace: "test-namespace",
			}

			// Lookup compatibility info
			compatibilityInfo, err := lookupImageAndFeatures(tc.pluginType, tc.clusterVersion)
			assert.NilError(t, err)

			// Verify SupportsTLSProfile matches our expectation
			assert.Equal(t, compatibilityInfo.SupportsTLSProfile, tc.supportsTLSProfile)

			// Create plugin info (this is where TLS profile gets applied)
			pluginInfo, err := PluginInfoBuilder(ctx, nil, nil, plugin, pluginConf, compatibilityInfo, tc.clusterVersion, logger)
			assert.NilError(t, err)
			assert.Assert(t, pluginInfo != nil)

			// Verify TLS configuration was applied correctly
			if tc.expectedTLSArgs {
				assert.Equal(t, pluginInfo.TLSMinVersion, tc.expectedMinVersion)
				assert.DeepEqual(t, pluginInfo.TLSCiphers, tc.expectedCiphers)
			} else {
				assert.Equal(t, pluginInfo.TLSMinVersion, "")
				assert.Assert(t, len(pluginInfo.TLSCiphers) == 0)
			}

			// Test deployment creation with TLS arguments
			deployment := newDeployment(*pluginInfo, "test-namespace", nil)
			args := deployment.Spec.Template.Spec.Containers[0].Args

			if tc.expectedTLSArgs && tc.expectedMinVersion != "" {
				assert.Assert(t, containsArg(args, "-tls-min-version="+tc.expectedMinVersion))
			} else {
				assert.Assert(t, !containsArgPrefix(args, "-tls-min-version="))
			}

			if tc.expectedTLSArgs && len(tc.expectedCiphers) > 0 {
				assert.Assert(t, containsArgPrefix(args, "-tls-cipher-suites="))
			} else {
				assert.Assert(t, !containsArgPrefix(args, "-tls-cipher-suites="))
			}

			// Base cert/key args should always be present
			assert.Assert(t, containsArg(args, "-cert=/var/serving-cert/tls.crt"))
			assert.Assert(t, containsArg(args, "-key=/var/serving-cert/tls.key"))
		})
	}
}