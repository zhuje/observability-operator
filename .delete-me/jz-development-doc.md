export TAG="1.1.0-dev-0.6.0"
export REPO="quay.io/jezhu/observability-operator"

## 0. Login into a cluster

# 1. Build image and push to repo 
make generate && \
make bundle && \
make operator-image bundle-image operator-push bundle-push  \
    IMAGE_BASE="quay.io/jezhu/observability-operator" \
    VERSION=1.1.0-dev-0.32.0

# 2. If logged into a cluster, deploy to cluster 
oc delete catalogsource observability-operator-catalog -n openshift-operators && \
operator-sdk cleanup observability-operator -n openshift-operators && \
operator-sdk run bundle \
    quay.io/jezhu/observability-operator-bundle:1.1.0-dev-0.32.0 \
    --install-mode AllNamespaces \
    --namespace openshift-operators \
    --security-context-config restricted


# 3. Apply UIPlugin

# (checked) throw error in logs 
oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
  namespace:
spec:
  type: Monitoring
  monitoring:
EOF


oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    perses:
      name: ""
      namespace: ""
EOF

oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    perses:
      name: "monitoring-plugin"
      namespace: "openshift-monitoring"
EOF

oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    perses:
      name: "perses-api-http""
      namespace: "perses-operator"
EOF



oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    alertmanager: 
      url: 'https://alertmanager.open-cluster-management-observability.svc:9095'
    thanosQuerier:
      url: 'https://rbac-query-proxy.open-cluster-management-observability.svc:8443'
    perses:
      name: "perses-api-http"
      namespace: "perses-operator"  
EOF

oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    alertmanager: 
      url: 'https://alertmanager.open-cluster-management-observability.svc:9095'
    thanosQuerier:
      url: 'https://rbac-query-proxy.open-cluster-management-observability.svc:8443'
    perses:
      name: ""
      namespace: ""
EOF

oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    alertmanager:
      url: 'https://alertmanager.open-cluster-management-observability.svc:9095'
    thanosQuerier:
      url: 'https://rbac-query-proxy.open-cluster-management-observability.svc:8443'
EOF



## Uninstall 

operator-sdk cleanup observability-operator -n openshift-operators
oc delete catalogsource observability-operator-catalog -n openshift-operators

## If uninstall is hung 
1. oc edit crd uiplugins.observability.openshift.io
n the editor, find the finalizers field under metadata, and remove any finalizers (it will look something like this):
```
metadata:
  finalizers:
  - kubernetes
```
After removing the finalizer(s), save and exit the editor. This should allow the CRD to be deleted.

2. Go the the UI > Installed Operator > manually delete the operator 


// +kubebuilder:validation:XValidation:rule="self.alertmanager != null && self.thanosQuerier != null || self.perses != null || (self.alertmanager != null && self.thanosQuerier != null && self.perses != null)",message="Either 'alertmanager' and 'thanosQuerier' are required, or 'perses' is required, or all three are required"


## If changing types.go / CRDs
This won't work in debug mode, you won't be able to see changes. You need to 
1. make generate 
2. make bundle 
3. rebuild and deploy. 

## To Update ConsolePlugin and monitoring-console-plugin deployment 
You need to deplete the monitoring UIPlugin and oc apply -f again to trigger the update 

## OU-571 branches 
perses-flag-dev-2
OU-571-perses-feature-flag-pr