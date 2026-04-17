# Manual Testing Guide: TLS Profile Refactor

## Prerequisites
- OpenShift cluster with observability-operator deployed
- UI plugins enabled (`openshift.enabled=true`)

## Test Scenarios

### 1. Verify TLS Profile Injection on Startup
```bash
# Check operator logs for TLS profile fetching
oc logs -n openshift-operators deployment/observability-operator | grep -i tls

# Expected log entries:
# - "fetched initial TLS profile" with minVersion and cipher count
# - No errors related to TLS profile fetching
```

### 2. Create Monitoring Plugin and Verify TLS Args
```bash
# Create a monitoring plugin
cat <<EOF | oc apply -f -
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: test-monitoring-tls
spec:
  type: monitoring
  monitoring:
    acm:
      enabled: true
      alertmanager:
        url: https://test-alertmanager:9093
      thanosQuerier:
        url: https://test-thanos:9091
EOF

# Wait for deployment creation
oc wait --for=condition=available deployment/test-monitoring-tls -n openshift-operators --timeout=300s

# Check deployment args for TLS configuration
oc get deployment/test-monitoring-tls -n openshift-operators -o yaml | grep -A 20 args:

# Expected for v4.19+ clusters:
# - "-tls-min-version=VersionTLS12" (or current cluster TLS setting)
# - "-tls-cipher-suites=<comma-separated-ciphers>"
# - "-cert=/var/serving-cert/tls.crt"
# - "-key=/var/serving-cert/tls.key"
```

### 3. Test TLS Profile Changes (Dynamic Update)
```bash
# Get current TLS profile
oc get apiserver cluster -o jsonpath='{.spec.tlsSecurityProfile}'

# Change TLS profile (if permitted)
oc patch apiserver cluster --type='merge' --patch='{"spec":{"tlsSecurityProfile":{"type":"Modern"}}}'

# Monitor operator logs - should see graceful restart
oc logs -n openshift-operators deployment/observability-operator -f | grep -i "tls\|restart\|graceful"

# Expected behavior:
# - "TLS security profile changed, triggering graceful restart"
# - Operator pod restart
# - New deployment with updated TLS settings
```

### 4. Test Non-Supporting Plugin (Expected No TLS Args)
```bash
# Create dashboards plugin (doesn't support TLS yet)
cat <<EOF | oc apply -f -
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: test-dashboards-no-tls
spec:
  type: dashboards
EOF

# Check deployment - should NOT have TLS args
oc get deployment/test-dashboards-no-tls -n openshift-operators -o yaml | grep -A 20 args:

# Expected:
# - NO "-tls-min-version" or "-tls-cipher-suites" args
# - Still has base "-cert" and "-key" args for HTTPS serving
```

### 5. Verify Compatibility Matrix
```bash
# Check operator logs for TLS profile application messages
oc logs -n openshift-operators deployment/observability-operator | grep "TLS profile not applied"

# Expected for non-supporting plugins:
# - "TLS profile not applied: plugin image does not support TLS profile flags"
```

## Validation Checklist

- [ ] Operator starts successfully with TLS profile fetch
- [ ] Monitoring plugin v4.19+ gets TLS arguments in deployment
- [ ] Monitoring plugin v4.18 does NOT get TLS arguments
- [ ] Other plugin types (dashboards, logging, etc.) do NOT get TLS arguments
- [ ] TLS profile changes trigger graceful operator restart
- [ ] New deployments after restart use updated TLS settings
- [ ] No regression in existing functionality

## Troubleshooting

### TLS Profile Fetch Fails
```bash
# Check RBAC permissions
oc get clusterrole observability-operator -o yaml | grep -A 5 "apiVersion.*config.openshift.io"

# Should see:
# - apiservers: [get, list, watch]
```

### Plugin Deployment Missing TLS Args
```bash
# Check compatibility matrix in operator code
# Only monitoring plugin v4.19+ should have SupportsTLSProfile: true

# Check cluster version
oc get clusterversion version -o jsonpath='{.status.desired.version}'
```

### Configuration Injection Issues
```bash
# Verify operator configuration
oc logs -n openshift-operators deployment/observability-operator | grep -i "configuration\|injection"

# Check for any WITH* function errors
```