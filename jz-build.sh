#!/bin/bash

# Define ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
ENDCOLOR='\033[0m' 

# Get the current date and time in 'monDD-HHMM' lowercase format
# For example: sep12-1216
TIMESTAMP=$(date +'%b%d-%H%M' | tr '[:upper:]' '[:lower:]')

VERSION="${VERSION:-1.0.0-dev-${TIMESTAMP}}"
IMG_BASE="${IMG_BASE:-"quay.io/jezhu/observability-operator"}"
IMAGE="${IMG_BASE}:${VERSION}"

print_title() {
  echo -e "\n${GREEN} =============================================== ${ENDCOLOR}\n"
  echo -e "${GREEN} $1 ${ENDCOLOR}"
  echo -e "\n${GREEN} =============================================== ${ENDCOLOR}\n"
}

# Build Bundle
print_title "Build Bundle: make operator-image bundle-image operator-push bundle-push"
make operator-image bundle-image operator-push bundle-push \
  IMG_BASE="${IMG_BASE}" \
  VERSION="${VERSION}"

# # edit ClusterServiceVersion to add "- -openshift.enabled=true"
# echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
# echo -e "${GREEN} Add - -openshift.enabled=true to ServiceClusterVersion  ${ENDCOLOR}"
# perl -i -0777 -pe '
# $exists = 1 if /--openshift.enabled=true/;
# if (!$exists) {
#     s/^(\s*)(- --namespace=\$\(NAMESPACE\).*)$/$1$2\n$1- --openshift.enabled=true/m;
# }
# END { 
#     print "exists = $exists\n"; 
#     print "Added --openshift.enabled=true\n" unless $exists; 
# }
# ' bundle/manifests/observability-operator.clusterserviceversion.yaml


# Delete Previous CatalogSource, Subscription, and ClusterServiceVersion
print_title "Delete Previous ClusterServiceVersion and Subscription"
# oc project openshift-operators
CAT_NAME=$(oc get catalogsource | grep 'observability-operator' | awk '{print $1}') && oc delete catalogsource ${CAT_NAME}
SUB_NAME=$(oc get subscriptions | grep 'observability-operator' | awk '{print $1}') && oc delete subscriptions ${SUB_NAME}
CSV_NAME=$(oc get clusterserviceversion | grep 'observability-operator' | awk '{print $1}') && oc delete clusterserviceversion ${CSV_NAME}

# OR Delete the whole operator 
operator-sdk cleanup observability-operator -n openshift-operators

# Run the bundle using the fully qualified image tag.
print_title "Run Bundle: operator-sdk run bundle" 
operator-sdk run bundle \
  ${IMG_BASE}-bundle:${VERSION} \
  --install-mode AllNamespaces \
  --namespace openshift-operators \
  --security-context-config restricted

# Edit ClusterServiceVersion to remove "- -openshift.enabled=true"
# echo -e "\n${GREEN} =============================================== ${ENDCOLOR}"
# echo -e "${GREEN} Remove: - -openshift.enabled=true from ServiceClusterVersion ${ENDCOLOR}"
# perl -i -pe 's/^\s*- --openshift.enabled=true\s*\n//m;' bundle/manifests/observability-operator.clusterserviceversion.yaml
