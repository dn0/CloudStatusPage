#!/usr/bin/env bash

set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")/.." || exit 64 ; pwd -P)"
APP="${1:-""}"
ENV="${2:-""}"
APP_DIR="${BASE_DIR}/k8s/${APP}"
declare -a HELM_PARAMS=(
    "--namespace" "default"
    "--install"
    "--cleanup-on-fail"
    "--wait"
    "--atomic"
    "--timeout" "60s"
)

if [[ -z "$ENV" ]] || [[ ! -f "${APP_DIR}/Chart.yaml" ]]; then
    echo "Usage: $0 <app> <environment>"
    exit 1
fi

if [[ -z "${APP_VERSION:-""}" ]]; then
    echo "APP_VERSION is not set" >&2
    exit 1
fi

if [[ -f "${APP_DIR}/secrets.yaml" ]]; then
    SOPS="sops"
    if ! which sops > /dev/null; then
        echo "Downloading sops"
        curl -f -s -S -L -o /tmp/sops https://github.com/getsops/sops/releases/download/v3.9.1/sops-v3.9.1.linux.amd64
        chmod 755 /tmp/sops
        SOPS="/tmp/sops"
    fi
    "$SOPS" decrypt --output "${APP_DIR}/secrets.dec.yaml" "${APP_DIR}/secrets.yaml"
    HELM_PARAMS+=("--values")
    HELM_PARAMS+=("${APP_DIR}/secrets.dec.yaml")
fi

BACKUP="/tmp/helm-deploy-${APP}-${APP_VERSION}-${ENV}-chart-yaml.$$"
cp -a "${APP_DIR}/Chart.yaml" "$BACKUP"
cleanup() {
    mv "$BACKUP" "${APP_DIR}/Chart.yaml"
    rm -f "${APP_DIR}/secrets.dec.yaml"
}
trap "cleanup" EXIT

sed -i "s/^appVersion: .*$/appVersion: \"${APP_VERSION}\"/g" "${APP_DIR}/Chart.yaml"

helm upgrade "$APP" "$APP_DIR" "${HELM_PARAMS[@]}" --set-string "image.tag=${APP_VERSION}"
