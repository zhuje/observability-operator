package uiplugin

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

var namespace = "openshift-operators"
var name = "monitoring"
var image = "quay.io/monitoring-foo-test:123"

var pluginConfigAll = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			ACM: uiv1alpha1.AdvancedClusterManagementReference{
				Enabled: true,
				Alertmanager: uiv1alpha1.AlertmanagerReference{
					Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
				},
				ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
					Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
				},
			},
			Perses: uiv1alpha1.PersesReference{
				Enabled:     true,
				ServiceName: "perses-api-http",
				Namespace:   "perses-operator",
			},
			Incidents: uiv1alpha1.IncidentsReference{
				Enabled: true,
			},
		},
	},
}

var pluginConfigPerses = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				Enabled:     true,
				ServiceName: "perses-api-http",
				Namespace:   "perses-operator",
			},
		},
	},
}

var pluginConfigPersesDefault = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				Enabled: true,
			},
		},
	},
}

var pluginConfigPersesNamespace = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				Namespace: "perses-operator",
			},
		},
	},
}

var pluginConfigPersesServiceName = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				ServiceName: "perses-api-http",
			},
		},
	},
}

var pluginConfigPersesNameSpace = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Perses: uiv1alpha1.PersesReference{
				ServiceName: "perses-api-http",
			},
		},
	},
}

var pluginConfigACM = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			ACM: uiv1alpha1.AdvancedClusterManagementReference{
				Enabled: true,
				Alertmanager: uiv1alpha1.AlertmanagerReference{
					Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
				},
				ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
					Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
				},
			},
		},
	},
}

var pluginConfigThanos = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			ACM: uiv1alpha1.AdvancedClusterManagementReference{
				ThanosQuerier: uiv1alpha1.ThanosQuerierReference{
					Url: "https://rbac-query-proxy.open-cluster-management-observability.svc:8443",
				},
			},
		},
	},
}

var pluginConfigAlertmanager = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			ACM: uiv1alpha1.AdvancedClusterManagementReference{
				Alertmanager: uiv1alpha1.AlertmanagerReference{
					Url: "https://alertmanager.open-cluster-management-observability.svc:9095",
				},
			},
		},
	},
}

var pluginConfigIncidents = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type: "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{
			Incidents: uiv1alpha1.IncidentsReference{
				Enabled: true,
			},
		},
	},
}

var pluginMalformed = &uiv1alpha1.UIPlugin{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "observability.openshift.io/v1alpha1",
		Kind:       "UIPlugin",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "monitoring-plugin",
	},
	Spec: uiv1alpha1.UIPluginSpec{
		Type:       "monitoring",
		Monitoring: &uiv1alpha1.MonitoringConfig{},
	},
}

func containsFeatureFlag(pluginInfo *UIPluginInfo) (bool, bool, bool) {
	acmAlertingFound, persesFound, incidentsFound := false, false, false
	var featuresIndex int

	// Loop through the array to find the index of "-features"
	for i, arg := range pluginInfo.ExtraArgs {
		if strings.Contains(arg, "-features") {
			fmt.Printf("Found '-features' at index: %d\n", i)
			featuresIndex = i
			break
		}
	}

	// Get "-features=" list from ExtraArgs field
	// (e.g. "-features='acm-alerting', 'perses-dashboards', 'incidents'")
	re := regexp.MustCompile(`-features=([a-zA-Z0-9,\-]+)`)
	featuresList := re.FindString(pluginInfo.ExtraArgs[featuresIndex])

	// Get individual feature strings, by spliting string after "=" and between ","
	features := strings.Split(strings.Split(featuresList, "=")[1], ",")

	// Check if features are listed
	for _, feature := range features {
		if feature == "acm-alerting" {
			acmAlertingFound = true
		}
		if feature == "perses-dashboards" {
			persesFound = true
		}
		if feature == "incidents" {
			incidentsFound = true
		}
	}

	return acmAlertingFound, persesFound, incidentsFound
}

func containsProxy(pluginInfo *UIPluginInfo) (bool, bool, bool) {
	alertmanagerFound, thanosFound, persesFound := false, false, false

	for _, proxy := range pluginInfo.Proxies {
		if proxy.Alias == "alertmanager-proxy" {
			alertmanagerFound = true
		}
		if proxy.Alias == "thanos-proxy" {
			thanosFound = true
		}
		if proxy.Alias == "perses" {
			persesFound = true
		}
	}
	return alertmanagerFound, thanosFound, persesFound
}

var acmVersion = "v2.11"
var features = []string{}
var clusterVersion = "v4.18"

func getPluginInfo(plugin *uiv1alpha1.UIPlugin, features []string) (*UIPluginInfo, error) {
	return createMonitoringPluginInfo(plugin, namespace, name, image, features, acmVersion, clusterVersion)
}

func TestCreateMonitoringPluginInfo(t *testing.T) {
	t.Run("Test createMonitoringPluginInfo with all monitoring configurations", func(t *testing.T) {
		pluginInfo, error := getPluginInfo(pluginConfigAll, features)
		assert.Assert(t, error == nil)

		fmt.Println("pluginInfo: ", pluginInfo)

		alertmanagerProxyFound, thanosProxyFound, persesProxyFound := containsProxy(pluginInfo)
		assert.Assert(t, alertmanagerProxyFound == true)
		assert.Assert(t, thanosProxyFound == true)
		assert.Assert(t, persesProxyFound == true)

		acmAlertingFlagFound, persesFlagFound, incidentsFlagFound := containsFeatureFlag(pluginInfo)
		assert.Assert(t, acmAlertingFlagFound == true)
		assert.Assert(t, persesFlagFound == true)
		assert.Assert(t, incidentsFlagFound == true)

	})

	t.Run("Test createMonitoringPluginInfo with AMC configuration only", func(t *testing.T) {
		pluginInfo, error := getPluginInfo(pluginConfigACM, features)
		assert.Assert(t, error == nil)

		alertmanagerProxyFound, thanosProxyFound, persesProxyFound := containsProxy(pluginInfo)
		assert.Assert(t, alertmanagerProxyFound == true)
		assert.Assert(t, thanosProxyFound == true)
		assert.Assert(t, persesProxyFound == false)

		acmAlertingFlagFound, persesFlagFound, incidentsFlagFound := containsFeatureFlag(pluginInfo)
		assert.Assert(t, acmAlertingFlagFound == true)
		assert.Assert(t, persesFlagFound == false)
		assert.Assert(t, incidentsFlagFound == false)
	})

	t.Run("Test createMonitoringPluginInfo with Perses configuration only", func(t *testing.T) {
		pluginInfo, error := getPluginInfo(pluginConfigPerses, features)
		assert.Assert(t, error == nil)

		alertmanagerProxyFound, thanosProxyFound, persesProxyFound := containsProxy(pluginInfo)
		assert.Assert(t, alertmanagerProxyFound == false)
		assert.Assert(t, thanosProxyFound == false)
		assert.Assert(t, persesProxyFound == true)

		acmAlertingFlagFound, persesFlagFound, incidentsFlagFound := containsFeatureFlag(pluginInfo)
		assert.Assert(t, acmAlertingFlagFound == false)
		assert.Assert(t, persesFlagFound == true)
		assert.Assert(t, incidentsFlagFound == false)

	})

	t.Run("Test createMonitoringPluginInfo with Incidents configuration only", func(t *testing.T) {
		pluginInfo, error := getPluginInfo(pluginConfigIncidents, features)
		fmt.Println("pluginInfo: ", pluginInfo)
		assert.Assert(t, error == nil)

		alertmanagerProxyFound, thanosProxyFound, persesProxyFound := containsProxy(pluginInfo)
		assert.Assert(t, alertmanagerProxyFound == false)
		assert.Assert(t, thanosProxyFound == false)
		assert.Assert(t, persesProxyFound == false)

		acmAlertingFlagFound, persesFlagFound, incidentsFlagFound := containsFeatureFlag(pluginInfo)
		assert.Assert(t, acmAlertingFlagFound == false)
		assert.Assert(t, persesFlagFound == false)
		assert.Assert(t, incidentsFlagFound == true)
	})

	/// HERE

	t.Run("Test createMonitoringPluginInfo with missing URL from thanos", func(t *testing.T) {
		errorMessage := AcmErrorMsg + PersesErrorMsg + ThanosEmptyMsg + PersesNameEmptyMsg + PersesNamespaceEmptyMsg

		// this should throw an error because thanosQuerier.URL is not set
		pluginInfo, error := getPluginInfo(pluginConfigAlertmanager, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), errorMessage)
	})

	t.Run("Test createMonitoringPluginInfo with missing URL from alertmanager ", func(t *testing.T) {
		errorMessage := AcmErrorMsg + PersesErrorMsg + AlertmanagerEmptyMsg + PersesNameEmptyMsg + PersesNamespaceEmptyMsg

		// this should throw an error because alertManager.URL is not set
		pluginInfo, error := getPluginInfo(pluginConfigThanos, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), errorMessage)
	})

	t.Run("Test createMonitoringPluginInfo with missing persesName ", func(t *testing.T) {
		errorMessage := AcmErrorMsg + PersesErrorMsg + AlertmanagerEmptyMsg + ThanosEmptyMsg + PersesNamespaceEmptyMsg

		// this should throw an error because persesName is not set
		pluginInfo, error := getPluginInfo(pluginConfigPersesNameSpace, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), errorMessage)
	})

	t.Run("Test createMonitoringPluginInfo with missing persesNamespace ", func(t *testing.T) {
		errorMessage := AcmErrorMsg + PersesErrorMsg + AlertmanagerEmptyMsg + ThanosEmptyMsg + PersesNameEmptyMsg

		// this should throw an error because persesNamespace is not set
		pluginInfo, error := getPluginInfo(pluginConfigPersesNamespace, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), errorMessage)
	})

	t.Run("Test createMonitoringPluginInfo with malform UIPlugin custom resource", func(t *testing.T) {
		errorMessage := AcmErrorMsg + PersesErrorMsg + AlertmanagerEmptyMsg + ThanosEmptyMsg + PersesNameEmptyMsg + PersesNamespaceEmptyMsg

		// this should throw an error because UIPlugin doesn't include alertmanager, thanos, or perses
		pluginInfo, error := getPluginInfo(pluginMalformed, features)
		assert.Assert(t, pluginInfo == nil)
		assert.Assert(t, error != nil)
		assert.Equal(t, error.Error(), errorMessage)
	})
}
