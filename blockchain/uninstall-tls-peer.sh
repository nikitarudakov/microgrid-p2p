#!/bin/bash

# === Set Namespace from Argument ===
NAMESPACE=${1:-org1}  # Default to 'org1' if no argument is provided
echo "🧹 Cleaning up peer org in namespace: $NAMESPACE"

# === Navigate to the Helm chart directory ===
cd ./blockchain/fabric-k8s/ || exit 1

# === Uninstall Helm Releases ===
echo "🗑️ Uninstalling Helm releases in $NAMESPACE..."

helm uninstall "$NAMESPACE-peer0"       || echo "⚠️ Could not uninstall $NAMESPACE-peer0"
helm uninstall "$NAMESPACE-org-ca"      || echo "⚠️ Could not uninstall $NAMESPACE-org-ca"
helm uninstall "$NAMESPACE-tls-ca"      || echo "⚠️ Could not uninstall $NAMESPACE-tls-ca"
helm uninstall "$NAMESPACE-persistence" || echo "⚠️ Could not uninstall $NAMESPACE-persistence"

# === Wait a moment for cleanup ===
sleep 3

# === Delete Namespace ===
echo "🧨 Deleting namespace $NAMESPACE..."
kubectl delete namespace "$NAMESPACE" --wait=true || echo "⚠️ Failed to delete namespace $NAMESPACE"

echo "✅ Peer org cleanup complete."