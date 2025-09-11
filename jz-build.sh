#!/bin/bash

# Define ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
ENDCOLOR='\033[0m' 

# Define variables with default values.
# The ':-' syntax provides a default if the variable is not set.
VERSION="${VERSION:-1.0.0-dev-sept10}"
IMG_BASE="${IMG_BASE:-"quay.io/jezhu/observability-operator"}"
IMAGE="${IMG_BASE}:${VERSION}"

# Delete Previous CatalogSource and Subscription 
echo -e "${GREEN} Delete Previous ClusterServiceVersion and Subscription  ${ENDCOLOR}"
oc project openshift-operators
CSV_NAME=$(oc get catalogsource | grep 'observability-operator' | awk '{print $1}') && oc delete catalogsource ${CSV_NAME}
SUB_NAME=$(oc get subscriptions | grep 'observability-operator' | awk '{print $1}') && oc delete subscriptions ${SUB_NAME}



# Build Bundle
echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
echo -e "${GREEN} Build Bundle: make operator-image bundle-image operator-push bundle-push ${ENDCOLOR}"
make operator-image bundle-image operator-push bundle-push \
  IMG_BASE="${IMG_BASE}" \
  VERSION="${VERSION}"


# edit ClusterServiceVersion to add "- -openshift.enabled=true"
echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
echo -e "${GREEN} Add - -openshift.enabled=true to ServiceClusterVersion  ${ENDCOLOR}"
perl -i -0777 -pe '
$exists = 1 if /--openshift.enabled=true/;
if (!$exists) {
    s/^(\s*)(- --namespace=\$\(NAMESPACE\).*)$/$1$2\n$1- --openshift.enabled=true/m;
}
END { 
    print "exists = $exists\n"; 
    print "Added --openshift.enabled=true\n" unless $exists; 
}
' bundle/manifests/observability-operator.clusterserviceversion.yaml

# Run the bundle using the fully qualified image tag.
echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
echo -e "${GREEN} Run Bundle: operator-sdk run bundle ${ENDCOLOR}" 
operator-sdk run bundle \
  "${IMG_BASE}-bundle:${VERSION}" \
  --install-mode AllNamespaces \
  --namespace openshift-operators \
  --security-context-config restricted

# Edit ClusterServiceVersion to remove "- -openshift.enabled=true"
# echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
# echo -e "${GREEN} Remove: - -openshift.enabled=true from ServiceClusterVersion ${ENDCOLOR}"
# perl -i -pe 's/^\s*- --openshift.enabled=true\s*\n//m;' bundle/manifests/observability-operator.clusterserviceversion.yaml
