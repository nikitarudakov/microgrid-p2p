#!/bin/bash

# Exit immediately on any error
set -e

# === Set Namespace/Org from Argument ===
NAMESPACE=${1:-org1}
echo "📦 Using namespace: $NAMESPACE"

# === Navigate to Helm Chart Directory ===
cd ./blockchain/fabric-k8s/ || exit 1

# === Create Namespace ===
echo "📦 Creating namespace $NAMESPACE"
kubectl create namespace "$NAMESPACE" || true

# === Deploy Persistence PVC ===
echo "📁 Installing persistence chart"
helm install "$NAMESPACE-persistence" ./ca/charts/persistence --set namespace="$NAMESPACE"

# === Wait for PVC to be Bound ===
echo "⏳ Waiting for PVC to be bound..."
PVC_NAME=$(kubectl get pvc -n "$NAMESPACE" -o jsonpath='{.items[0].metadata.name}')
kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/"$PVC_NAME" -n "$NAMESPACE" --timeout=60s

# === Deploy TLS CA ===
echo "🔐 Installing TLS CA"
helm install "$NAMESPACE-tls-ca" ./ca/charts/tls \
  --set namespace="$NAMESPACE" \
  --set ca.name="tls-${NAMESPACE}-ca" \
  --set node.type=peer \
  --set node.name=peer0 \
  --set node.secret=peer0pw

# === Wait for TLS CA Enrollment Job ===
echo "⏳ Waiting for 'tls-ca-enrollment' Job to complete in $NAMESPACE..."
kubectl wait --for=condition=complete job/tls-ca-enrollment -n "$NAMESPACE" --timeout=60s

# === Deploy Org CA ===
echo "🔐 Installing Org CA"
helm install "$NAMESPACE-org-ca" ./ca/charts/org \
  --set namespace="$NAMESPACE" \
  --set ca.name="${NAMESPACE}-ca" \
  --set node.type=peer \
  --set node.name=peer0 \
  --set node.secret=peer0pw

# === Wait for Org CA Enrollment Job ===
echo "⏳ Waiting for 'org-ca-enrollment' Job to complete in $NAMESPACE..."
kubectl wait --for=condition=complete job/org-ca-enrollment -n "$NAMESPACE" --timeout=60s

# === Deploy Peer ===
echo "📦 Deploying Peer"
helm install "$NAMESPACE-peer0" ./peer \
  --set namespace="$NAMESPACE" \
  --set msp.id="${NAMESPACE^}MSP" \
  --set org.address=${NAMESPACE}.svc.cluster.local