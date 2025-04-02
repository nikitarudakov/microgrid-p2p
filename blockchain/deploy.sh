#!/bin/bash

cd ./blockchain/fabric-k8s/

# Deploy Root CA
echo "Deploy root"
helm install root-ca ./root

# Deploy Org1CA, register, enroll its identities
kubectl create namespace org1 || true

echo "Deploy Org1"
helm install org1-ca ./base --set namespace=org1 --set name='Org1 CA' --create-namespace

kubectl create namespace org2 || true

echo "Deploy Org2"
helm install org2-ca ./base --set namespace=org2 --set name='Org2 CA' --create-namespace

echo "Deploy Org1 Peer"
helm install peer0-org1 ./peer --set namespace=org1 --set peer.org=org1 --set peer.mspID=Org1MSP

echo "Deploy Org2 Peer"
helm install peer0-org2 ./peer --set namespace=org2 --set peer.org=org2 --set peer.mspID=Org2MSP