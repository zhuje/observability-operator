package operator

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"gotest.tools/v3/assert"
)

func TestWithTLSProfile(t *testing.T) {
	testCases := []struct {
		name       string
		tlsProfile configv1.TLSProfileSpec
		expected   configv1.TLSProfileSpec
	}{
		{
			name: "TLS profile with min version and ciphers",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS12,
				Ciphers:       []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"},
			},
			expected: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS12,
				Ciphers:       []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"},
			},
		},
		{
			name: "TLS profile with TLS 1.3 only",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS13,
				Ciphers:       []string{},
			},
			expected: configv1.TLSProfileSpec{
				MinTLSVersion: configv1.VersionTLS13,
				Ciphers:       []string{},
			},
		},
		{
			name: "Empty TLS profile",
			tlsProfile: configv1.TLSProfileSpec{
				MinTLSVersion: "",
				Ciphers:       nil,
			},
			expected: configv1.TLSProfileSpec{
				MinTLSVersion: "",
				Ciphers:       nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a base configuration
			config := NewOperatorConfiguration()

			// Apply the TLS profile using WithTLSProfile function
			withTLSProfile := WithTLSProfile(tc.tlsProfile)
			withTLSProfile(config)

			// Verify the TLS profile was correctly applied to UIPlugins configuration
			assert.DeepEqual(t, config.UIPlugins.TLSProfile, tc.expected)
		})
	}
}

func TestOperatorConfigurationChaining(t *testing.T) {
	// Test that WithTLSProfile can be chained with other configuration functions
	tlsProfile := configv1.TLSProfileSpec{
		MinTLSVersion: configv1.VersionTLS12,
		Ciphers:       []string{"TLS_AES_128_GCM_SHA256"},
	}

	config := NewOperatorConfiguration(
		WithNamespace("test-namespace"),
		WithMetricsAddr(":8080"),
		WithTLSProfile(tlsProfile),
		WithFeatureGates(FeatureGates{
			OpenShift: OpenShiftFeatureGates{
				Enabled: true,
			},
		}),
	)

	// Verify all configurations were applied
	assert.Equal(t, config.Namespace, "test-namespace")
	assert.Equal(t, config.MetricsAddr, ":8080")
	assert.DeepEqual(t, config.UIPlugins.TLSProfile, tlsProfile)
	assert.Equal(t, config.FeatureGates.OpenShift.Enabled, true)
}