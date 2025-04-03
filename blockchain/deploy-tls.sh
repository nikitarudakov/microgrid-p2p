#!/bin/bash

# Navigate to the Helm chart directory
cd ./blockchain/fabric-k8s/

# === Deploy Orderer CA ===
echo "📦 Creating namespace for Orderer Org"
kubectl create namespace orderer-org || true

helm install tls-ca ./tls --set namespace=orderer-org