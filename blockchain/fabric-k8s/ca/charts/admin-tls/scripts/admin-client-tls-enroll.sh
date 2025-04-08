#!/bin/sh

clean_msp() {
  local msp_dir=$1
  echo "🧹 Cleaning MSP at $msp_dir"

  rm -rf ${msp_dir}/keystore/*
  rm -rf ${msp_dir}/signcerts/*
  rm -rf ${msp_dir}/tlscacerts/*
}

rm -rf admin-client

export CA_HOST="{{ .Values.name }}.${ORG_DOMAIN}:7054"
export NODE_MSPDIR="tls"

echo "⏳ Waiting for TLS CA to become ready..."
sleep 10

# Admin enrollment
clean_msp "${FABRIC_CA_CLIENT_HOME}/tlsca/msp"
fabric-ca-client enroll \
  -u "https://{{ .Values.admin.user }}:{{ .Values.admin.pass }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --enrollment.profile tls \
  --mspdir "tlsca/msp" \
  -d

# Register Admin Client identity
fabric-ca-client register \
  --id.name "{{ .Values.client.name }}" \
  --id.secret "{{ .Values.client.secret }}" \
  --id.type client \
  -u "https://${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "tlsca/msp" \
  -d

# Enroll Admin Client identity
clean_msp "${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}"
fabric-ca-client enroll \
  -u "https://{{ .Values.client.name }}:{{ .Values.client.secret }}@${CA_HOST}" \
  --enrollment.profile tls \
  --csr.hosts "admin-tools.{{ .Values.namespace }}.svc.cluster.local,localhost" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "${NODE_MSPDIR}" \
  -d

echo "⏳ Waiting for Nodes Private Key to become ready..."
sleep 10

# Rename node key and certs
KEY_PATH="${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}"
echo "⏳ Waiting for Node Private Key to appear..."
while [ -z "$(ls -A ${KEY_PATH}/keystore 2>/dev/null)" ]; do
  echo "🔁 Still waiting for private key in ${KEY_PATH}/keystore..."
  sleep 1
done

mv ${KEY_PATH}/keystore/*_sk "${FABRIC_CA_CLIENT_HOME}/client-tls-key.pem"
mv ${KEY_PATH}/signcerts/cert.pem "${FABRIC_CA_CLIENT_HOME}/client-tls-cert.pem"
cp "${FABRIC_CA_HOME}/ca-cert.pem" "${FABRIC_CA_CLIENT_HOME}/client-tls-ca-cert.pem"

echo "✅ TLS enrollment and registration complete"