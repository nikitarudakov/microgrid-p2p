#!/bin/bash

# Navigate to the Helm chart directory
cd ./blockchain/fabric-k8s/

# === Deploy TLS Orderer CA ===
echo "📦 Creating namespace for Orderer Org"
kubectl create namespace orderer-org || true

helm install tls-ca ./tls --set namespace=orderer-org --set ca.name=tls-orderer-ca

# === Jobs ===
echo "⏳ Waiting for 'tls-ca-enrollment' Job to complete in orderer-org..."
kubectl wait --for=condition=complete job/tls-ca-enrollment -n orderer-org --timeout=60s

sleep 5

# === Deploy Orderer Org CA ===
echo "🔐 Deploying Orderer Org CA"
helm install orderer-ca ./org-ca --set namespace=orderer-org --set ca.name=orderer-ca

# === Jobs ===
echo "⏳ Waiting for 'org-ca-enrollment' Job to complete in orderer-org..."
kubectl wait --for=condition=complete job/org-ca-enrollment -n orderer-org --timeout=60s