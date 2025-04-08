#!/bin/sh

clean_msp() {
  local msp_dir=$1
  echo "🧹 Cleaning MSP at $msp_dir"

  rm -rf ${msp_dir}/keystore/*
  rm -rf ${msp_dir}/signcerts/*
  rm -rf ${msp_dir}/tlscacerts/*
}

export CA_HOST="{{ .Values.name }}.${ORG_DOMAIN}:7054"
export NODE_MSPDIR="{{ .Values.node.type}}s/{{ .Values.node.name }}.${ORG_DOMAIN}/msp"

echo "⏳ Waiting for Org CA to become ready..."
sleep 10

echo "⏳ Enrolling Org Root CA Admin..."
clean_msp "${FABRIC_CA_CLIENT_HOME}/ca/{{ .Values.root.user }}/msp"
fabric-ca-client enroll \
  -u "https://{{ .Values.root.user }}:{{ .Values.root.pass }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "ca/{{ .Values.root.user }}/msp" \
  -d
echo "✅ Root CA Admin enrolled"

sleep 5

echo "📌 Registering Org Admin with Root CA..."
fabric-ca-client register \
  --id.name "{{ .Values.org.admin.user }}" \
  --id.secret "{{ .Values.org.admin.pass }}" \
  --id.type admin \
  --id.attrs "hf.Registrar.Roles=*,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,hf.AffiliationMgr=true" \
  -u "https://${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "ca/{{ .Values.root.user }}/msp" \
  -d
echo "✅ Org Admin registered"

echo "⏳ Enrolling Org Admin..."
clean_msp "${FABRIC_CA_CLIENT_HOME}/msp"
fabric-ca-client enroll \
  -u "https://{{ .Values.org.admin.user }}:{{ .Values.org.admin.pass }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "msp" \
  -d
echo "✅ Org Admin enrolled"

echo "📌 Registering node '{{ .Values.node.name }}'..."
fabric-ca-client register \
  --id.name "{{ .Values.node.name }}" \
  --id.secret "{{ .Values.node.secret }}" \
  --id.type "{{ .Values.node.type }}" \
  -u "https://${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "msp" \
  -d
echo "✅ Node '{{ .Values.node.name }}' registered with Org CA"

echo "⏳ Enrolling node '{{ .Values.node.name }}'..."
clean_msp "${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}"
fabric-ca-client enroll \
  -u "https://{{ .Values.node.name }}:{{ .Values.node.secret }}@${CA_HOST}" \
  --tls.certfiles "${FABRIC_CA_HOME}/ca-cert.pem" \
  --mspdir "${NODE_MSPDIR}" \
  -d
echo "✅ Node '{{ .Values.node.name }}' enrolled with Org CA"

echo "📁 Copying Node MSP config.yaml..."
cp /tmp/config/config.yaml "${FABRIC_CA_CLIENT_HOME}/${NODE_MSPDIR}/config.yaml"
cp /tmp/config/config.yaml "${FABRIC_CA_CLIENT_HOME}/msp/config.yaml"
echo "✅ MSP configuration copied"

echo "📁 Copying Root CA certificates..."
mkdir -p ${FABRIC_CA_CLIENT_HOME}/msp/cacerts
mkdir -p ${FABRIC_CA_CLIENT_HOME}/msp/tlscacerts
echo "✅ Root CA certificates copied"