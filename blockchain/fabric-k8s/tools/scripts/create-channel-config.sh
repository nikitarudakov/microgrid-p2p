#!/bin/sh

# Set orderer info
export ORDERER_DOMAIN=orderer.svc.cluster.local
export ORDERER_URL=orn1.orderer.svc.cluster.local
export ORDERER_ADMIN_PORT=9443

# Generate channel genesis block
configtxgen -profile SampleSingleMSPEtcd -outputBlock genesis_block.pb -channelID channel1

# Set TLS admin certs
OSN_TLS_CA_ROOT_CERT=./organizations/${ORDERER_DOMAIN}/orderers/${ORDERER_URL}/tls/ca.crt
export ADMIN_TLS_SIGN_CERT=./admin-client/client-tls-cert.pem
export ADMIN_TLS_PRIVATE_KEY=./admin-client/client-tls-key.pem

# Join orderer to the channel
osnadmin channel join \
  --channelID channel1 \
  --config-block genesis_block.pb \
  -o ${ORDERER_URL}:${ORDERER_ADMIN_PORT} \
  --ca-file "$OSN_TLS_CA_ROOT_CERT" \
  --client-cert "$ADMIN_TLS_SIGN_CERT" \
  --client-key "$ADMIN_TLS_PRIVATE_KEY"

# Set peer environment
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=./organizations/org1.svc.cluster.local/peers/peer0.org1.svc.cluster.local/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=./organizations/org1.svc.cluster.local/users/Admin@org1.svc.cluster.local/msp
export CORE_PEER_ADDRESS=peer0.org1.svc.cluster.local:7051

cp ./config/core.yaml ./core.yaml
export FABRIC_CFG_PATH=/etc/hyperledger/fabric

# Fetch channel block from orderer
peer channel fetch 0 \
  channel-artifacts/channel1.block \
  -o ${ORDERER_URL}:7050 \
  --ordererTLSHostnameOverride ${ORDERER_URL} \
  -c channel1 \
  --tls --cafile "${OSN_TLS_CA_ROOT_CERT}"
