#!/bin/bash

# Define variables with default values.
# The ':-' syntax provides a default if the variable is not set.
VERSION="${VERSION:-1.0.0-dev-sept10}"
IMG_BASE="${IMG_BASE:-"quay.io/jezhu/observability-operator"}"
IMAGE="${IMG_BASE}:${VERSION}"

# Print the image tag to be used.
echo "IMAGE: ${IMAGE}"

# Call the make command to build and push images.
# These variables are passed directly to the Makefile.
make operator-image bundle-image operator-push bundle-push \
  IMG_BASE="${IMG_BASE}" \
  VERSION="${VERSION}"


# edit ClusterServiceVersion to add "- -openshift.enabled=true"
echo "SCRIPT LOG: edit ClusterServiceVersion to add: - -openshift.enabled=true"
perl -i.bak -0777 -pe '
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
echo " =========================================== "
echo " Starting command .... operator-sdk run bundle" 
operator-sdk run bundle \
  "${IMAGE}" \
  --install-mode AllNamespaces \
  --namespace openshift-operators \
  --security-context-config restricted

# edit ClusterServiceVersion to remove "- -openshift.enabled=true"
echo "SCRIPT LOG: edit ClusterServiceVersion to remove: - -openshift.enabled=true "
perl -i -pe 's/^\s*- --openshift.enabled=true\s*\n//m;' bundle/manifests/observability-operator.clusterserviceversion.yaml

  