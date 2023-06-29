#!/bin/bash

mkdir -p ./.bin
export PATH="./.bin:$PATH"

set -euxo pipefail

# renovate: datasource=github-releases depName=kubeconform packageName=yannh/kubeconform
KUBECONFORM_VERSION=v0.6.2

CHART_DIR="./operations/helm/charts/fider"

SCHEMA_LOCATION="https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/"

# install kubeconform
curl --silent --show-error --fail --location --output /tmp/kubeconform.tar.gz https://github.com/yannh/kubeconform/releases/download/"${KUBECONFORM_VERSION}"/kubeconform-linux-amd64.tar.gz
tar -C .bin/ -xf /tmp/kubeconform.tar.gz kubeconform

# validate chart
(cd "${CHART_DIR}"; helm dependency build)
  helm template \
    --values "${CHART_DIR}/values.yaml" \
    "${CHART_DIR}" | kubeconform \
      --strict \
      --ignore-missing-schemas \
      --kubernetes-version "${KUBERNETES_VERSION#v}" \
      --schema-location "${SCHEMA_LOCATION}"
