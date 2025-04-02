#!/bin/bash

# Navigate to the Helm chart directory
cd ./blockchain/fabric-k8s/

# === Deploy Root CA ===
echo "🚀 Deploying Root CA"
helm install root-ca ./root

# === Deploy Org1 CA ===
echo "📦 Creating namespace for Org1"
kubectl create namespace org1 || true

echo "🔐 Deploying Org1 CA (Intermediate CA)"
helm install org1-ca ./base \
  --set namespace=org1 \
  --set name=org1-ca

# === Deploy Org2 CA ===
echo "📦 Creating namespace for Org2"
kubectl create namespace org2 || true

echo "🔐 Deploying Org2 CA (Intermediate CA)"
helm install org2-ca ./base \
  --set namespace=org2 \
  --set name=org2-ca

# === Deploy Peers ===

kubectl wait --for=condition=ready pod -l app=org1-ca -n org1 --timeout=60s

echo "🎯 Deploying Peer0 for Org1"
helm install peer0-org1 ./peer \
  --set namespace=org1 \
  --set peer.org=org1 \
  --set peer.mspID=Org1MSP \

kubectl wait --for=condition=ready pod -l app=org2-ca -n org2 --timeout=60s

echo "🎯 Deploying Peer0 for Org2"
helm install peer0-org2 ./peer \
  --set namespace=org2 \
  --set peer.org=org2 \
  --set peer.mspID=Org2MSP