#!/bin/bash

# === Set Namespace from Argument ===
NAMESPACE=${1:-orderer}  # Default to 'orderer' if no argument is provided
echo "🧹 Cleaning up namespace: $NAMESPACE"

# === Navigate to the Helm chart directory ===
cd ./blockchain/fabric-k8s/ || exit 1

# === Uninstall Helm Releases ===
echo "🗑️ Uninstalling Helm releases in $NAMESPACE..."

helm uninstall "$NAMESPACE-orderer0"    || echo "⚠️ Could not uninstall $NAMESPACE-orderer0"
helm uninstall "$NAMESPACE-orderer-ca"  || echo "⚠️ Could not uninstall $NAMESPACE-orderer-ca"
helm uninstall "$NAMESPACE-tls-ca"      || echo "⚠️ Could not uninstall $NAMESPACE-tls-ca"
helm uninstall "$NAMESPACE-persistence" || echo "⚠️ Could not uninstall $NAMESPACE-persistence"

# === Wait a moment for resources to be released ===
sleep 3

# === Delete Namespace ===
echo "🧨 Deleting namespace $NAMESPACE..."
kubectl delete namespace "$NAMESPACE" --wait=true || echo "⚠️ Failed to delete namespace $NAMESPACE"

echo "✅ Cleanup complete."