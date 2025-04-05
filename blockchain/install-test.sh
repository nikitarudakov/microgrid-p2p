cd ./blockchain/fabric-k8s/

minikube start

helm install nfs ./nfs

# === Deploy TLS Orderer CA ===
echo "📦 Creating namespace for Org1"
kubectl create namespace org1 || true

helm install org1-persistence ./ca/charts/persistence/

helm install org1-tls-ca ./ca/charts/tls/