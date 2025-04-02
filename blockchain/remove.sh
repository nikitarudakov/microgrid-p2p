#!/bin/bash

cd ./blockchain/fabric-k8s/

# Deploy Root CA
echo "Remove root"
helm uninstall root-ca

# Deploy Org1CA, register, enroll its identities
echo "Remove Org1"
helm uninstall org1-ca

#echo "Deploy Org2"
#helm install org2-ca ./base --set container.name=org2-ca --set name='Org2 CA' --set identities.affiliation=org2