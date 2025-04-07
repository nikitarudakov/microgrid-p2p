#!/bin/bash

# Exit immediately on any error
set -e

# === Set Namespace from Argument ===
NAMESPACE=${1:-orderer}  # Default to 'orderer' if no argument is provided
echo "📦 Using namespace: $NAMESPACE"

# === Navigate to the Helm chart directory ===
cd ./blockchain/fabric-k8s/ || exit 1

# === Create Namespace ===
echo "📦 Creating namespace $NAMESPACE"
kubectl create namespace "$NAMESPACE" || true

# === Deploy Persistence VolumeClaim ===
echo "📁 Installing persistence chart"
helm install $NAMESPACE-persistence ./ca/charts/persistence --set namespace="$NAMESPACE"

sleep 5

# === Wait for PVC to be Bound ===
echo "⏳ Waiting for PVC to be bound..."
PVC_NAME=$(kubectl get pvc -n "$NAMESPACE" -o jsonpath='{.items[0].metadata.name}')
kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/ca-pvc -n "$NAMESPACE" --timeout=90s

# === Deploy TLS Orderer CA ===
echo "🔐 Installing TLS Orderer CA"
helm install $NAMESPACE-tls-ca ./ca/charts/tls \
  --set namespace="$NAMESPACE" \
  --set ca.name=tls-orderer-ca \
  --set node.name=orn1

# === Wait for TLS CA Enrollment Job ===
echo "⏳ Waiting for 'tls-ca-enrollment' Job to complete in $NAMESPACE..."
kubectl wait --for=condition=complete job/tls-ca-enrollment -n "$NAMESPACE" --timeout=90s

# === Deploy Orderer Org CA ===
echo "🔐 Installing Orderer Org CA"
helm install $NAMESPACE-orderer-ca ./ca/charts/org \
  --set namespace="$NAMESPACE" \
  --set ca.name=orderer-ca \
  --set node.name=orn1

# === Wait for Org CA Enrollment Job ===
echo "⏳ Waiting for 'org-ca-enrollment' Job to complete in $NAMESPACE..."
kubectl wait --for=condition=complete job/org-ca-enrollment -n "$NAMESPACE" --timeout=90s

# === Deploy Orderer ===
echo "📦 Deploying Orderer"
helm install $NAMESPACE-orderer0 ./orderer \
  --set name=orn1 \
  --set namespace="$NAMESPACE" \
  --set ports.listen=7051 \
  --set ports.admin=9443 \
  --set org.address=$NAMESPACE.svc.cluster.local

#echo "📦 Deploying Orderer"
#helm install $NAMESPACE-orderer0 ./orderer \
#  --set name=orn2 \
#  --set namespace="$NAMESPACE" \
#  --set ports.listen=7051 \
#  --set ports.admin=9444 \
#  --set org.address=$NAMESPACE.svc.cluster.local
