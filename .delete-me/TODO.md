[x] 1. Make thanosQuerier and alertmanager URL required in CRD 
[x] 2. Create ‘enable’ attribute on acm and perses 
[x] 3. Add ‘incidents’ to the CR as well 
[x] 4. Change Perses ‘name’ to ‘serviceName’
[x] 5. If Perses is enabled and no serviceName/Namespace is set then use the default serviceName and Namespace (‘perses-api-http’, ‘perses’). Namespace: Perses, is where the backend where live . Perses-operator will be deployed in the same namespace as COO (openshift-operator) 
[x] 6. getConfigError() to getAcmConfigError() and getPersesConfigError(). 


### When an CRD error is thrown 
1. Make thanosQuerier and alertmanager URL required in CRD - works only if you leave "url:" blank 
❯  oc apply -f - <<EOF
apiVersion: observability.openshift.io/v1alpha1
kind: UIPlugin
metadata:
  name: monitoring
spec:
  type: Monitoring
  monitoring:
    alertmanager:
      url:
    thanosQuerier:
      url: 'https://rbac-query-proxy.open-cluster-management-observability.svc:8443'
EOF

The UIPlugin "monitoring" is invalid: 
* spec.monitoring.alertmanager.url: Required value
* <nil>: Invalid value: "null": some validation rules were not checked because the object was invalid; correct the existing errors to complete validation