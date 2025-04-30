make generate && \
make bundle && \
make operator-image bundle-image operator-push bundle-push  \
    IMAGE_BASE="quay.io/jezhu/observability-operator" \
    VERSION=1.1.0-otel-severity-0.0.1


oc delete catalogsource observability-operator-catalog -n openshift-operators && \
operator-sdk cleanup observability-operator -n openshift-operators && \
operator-sdk run bundle \
    quay.io/jezhu/observability-operator-bundle:1.1.0-otel-severity-0.0.1 \
    --install-mode AllNamespaces \
    --namespace openshift-operators \
    --security-context-config restricted

1. Check if missing upstream description with Kubebuilder 
/Users/jezhu/go/pkg/mod/github.com/perses/perses@v0.51.0-beta.0/pkg/model/api/config/dashboard.go:33:2: encountered struct field "jsonEval" without JSON tag in type "CustomLintRule"
/Users/jezhu/go/pkg/mod/github.com/perses/perses@v0.51.0-beta.0/pkg/model/api/config/dashboard.go:37:2: encountered struct field "celProgram" without JSON tag in type "CustomLintRule"
2. make perses-ops-crd creates deploy/perses/crds -- check these yamls 
3. make bundle creates deploy/crds/observability-operators/csv



# modification needed for development 
1. flag.BoolVar(&openShiftEnabled, "openshift.enabled", true, "Enable OpenShift specific features such as Console Plugins.")
2. You can pass in specific images with the flag -args in the ClusterServiceVersion YAML 
