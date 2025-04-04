#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Accept org name as the first argument, default to "org1" if not provided
ORG_NAME="${1:-org1}"

# Navigate to the Helm chart directory
cd ./blockchain/fabric-k8s/

# === Deploy TLS CA ===
echo "📦 Creating namespace for ${ORG_NAME}"
kubectl create namespace "${ORG_NAME}" || true

helm install tls-peer-${ORG_NAME}-ca ./tls-ca \
  --set namespace="${ORG_NAME}" \
  --set ca.name="tls-${ORG_NAME}-ca" \
  --set node.type=peer \
  --set node.name=peer0 \
  --set node.secret=peer0pw

# === Jobs ===
echo "⏳ Waiting for 'tls-ca-enrollment' Job to complete in ${ORG_NAME}..."
kubectl wait --for=condition=complete job/tls-ca-enrollment -n "${ORG_NAME}" --timeout=60s

sleep 5

# === Deploy Org CA ===
echo "🔐 Deploying Org CA"
helm install "${ORG_NAME}-ca" ./org-ca \
  --set namespace="${ORG_NAME}" \
  --set ca.name="${ORG_NAME}-ca" \
  --set node.type=peer \
  --set node.name=peer0 \
  --set node.secret=peer0pw

# === Jobs ===
echo "⏳ Waiting for 'org-ca-enrollment' Job to complete in ${ORG_NAME}..."
kubectl wait --for=condition=complete job/org-ca-enrollment -n "${ORG_NAME}" --timeout=60s

sleep 5

# === Deploy Peer ===
echo "🔐 Deploying Peer"
helm install ${ORG_NAME}-peer0 ./peer \
  --set namespace="${ORG_NAME}" \
  --set msp.id="peer0-${ORG_NAME}"