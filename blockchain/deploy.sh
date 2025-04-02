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

# === Jobs ===
echo "⏳ Waiting for 'register-enroll-identities' Job to complete in org1..."
kubectl wait --for=condition=complete job/register-enroll-identities -n org1 --timeout=60s

echo "⏳ Waiting for 'register-enroll-identities' Job to complete in org2..."
kubectl wait --for=condition=complete job/register-enroll-identities -n org2 --timeout=60s

# === Deploy Peers ===
echo "🎯 Deploying Peer0 for Org1"
helm install peer0-org1 ./peer \
  --set namespace=org1 \
  --set peer.org=org1 \
  --set peer.mspID=Org1MSP

echo "🎯 Deploying Peer0 for Org2"
helm install peer0-org2 ./peer \
  --set namespace=org2 \
  --set peer.org=org2 \
  --set peer.mspID=Org2MSP