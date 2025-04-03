#!/bin/bash

# Navigate to the Helm chart directory
cd ./blockchain/fabric-k8s/

# === Deploy Root CA ===
echo "🚀 Deploying Root CA"
helm install root-ca ./root

# === Deploy Orderer CA ===
echo "📦 Creating namespace for Orderer Org"
kubectl create namespace orderer-org || true

echo "🔐 Deploying Orderer Org CA (Intermediate CA)"
helm install orderer-ca ./base \
  --set namespace=orderer-org \
  --set name=orderer-ca \
  --set register.includePeer=false \
  --set register.includeOrderer=true