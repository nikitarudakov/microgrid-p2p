#!/bin/sh

clean_msp() {
  local msp_dir=$1
  echo "🧹 Cleaning MSP at $msp_dir"

  rm -rf ${msp_dir}/keystore/*
  rm -rf ${msp_dir}/signcerts/*
  rm -rf ${msp_dir}/tlscacerts/*
}

export CA_HOST="{{ .Values.name }}.${ORG_DOMAIN}:7054"
export NODE_MSPDIR="{{ .Values.node.type}}s/{{ .Values.node.name }}.${ORG_DOMAIN}/tls"

echo "⏳ Waiting for TLS CA to become ready..."
sleep 10

# Admin enrollment
clean_msp "${FABRIC_CA_CLIENT_HOME}/tlsca/{{ .Values.admin.user }}/msp"
fabric-ca-client enroll \
  -u "https://{{ .Values.admin.user }}:{{ .Values.admin.pass }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --enrollment.profile tls \
  --mspdir "tlsca/{{ .Values.admin.user }}/msp" \
  -d

# Register Org bootstrap identity
fabric-ca-client register \
  --id.name "{{ .Values.org.boot.name }}" \
  --id.secret "{{ .Values.org.boot.secret }}" \
  -u "https://${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "tlsca/{{ .Values.admin.user }}/msp" \
  -d

# Enroll Org bootstrap identity
clean_msp "${FABRIC_CA_CLIENT_HOME}/tlsca/{{ .Values.org.boot.name }}/msp"
fabric-ca-client enroll \
  -u "https://{{ .Values.org.boot.name }}:{{ .Values.org.boot.secret }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --enrollment.profile tls \
  --csr.hosts "*.{{ .Values.namespace }}.svc.cluster.local,localhost" \
  --mspdir "tlsca/{{ .Values.org.boot.name }}/msp" \
  -d

# Rename Org bootstrap key
KEY_PATH="${FABRIC_CA_CLIENT_HOME}/tlsca/{{ .Values.org.boot.name }}/msp/keystore"
mv ${KEY_PATH}/* "${KEY_PATH}/key.pem"

# Register node
fabric-ca-client register \
  --id.name "{{ .Values.node.name }}" \
  --id.secret "{{ .Values.node.secret }}" \
  --id.type "{{ .Values.node.type }}" \
  -u "https://${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "tlsca/{{ .Values.admin.user }}/msp" \
  -d

# Enroll node
clean_msp "${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}"
fabric-ca-client enroll \
  -u "https://{{ .Values.node.name }}:{{ .Values.node.secret }}@${CA_HOST}" \
  --enrollment.profile tls \
  --csr.hosts "{{ .Values.node.name }}.{{ .Values.namespace }}.svc.cluster.local,localhost" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "${NODE_MSPDIR}" \
  -d

# Rename node key and certs
KEY_PATH="${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}"
echo "⏳ Waiting for Node Private Key to appear..."
while [ -z "$(ls -A ${KEY_PATH}/keystore 2>/dev/null)" ]; do
  echo "🔁 Still waiting for private key in ${KEY_PATH}/keystore..."
  sleep 1
done

mv ${KEY_PATH}/keystore/*_sk "${KEY_PATH}/server.key"
mv ${KEY_PATH}/signcerts/cert.pem "${KEY_PATH}/server.crt"
cp "${FABRIC_CA_HOME}/ca-cert.pem" "${KEY_PATH}/ca.crt"

echo "✅ TLS enrollment and registration complete"