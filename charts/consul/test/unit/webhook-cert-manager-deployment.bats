#!/usr/bin/env bats

load _helpers

@test "webhookCertManager/Deployment: enabled by default" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      . | tee /dev/stderr |
      yq 'length > 0' | tee /dev/stderr)
  [ "${actual}" = "true" ]
}

@test "webhookCertManager/Deployment: enabled with controller.enabled=true and connectInject.enabled=false" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      . | tee /dev/stderr |
      yq 'length > 0' | tee /dev/stderr)
  [ "${actual}" = "true" ]
}

@test "webhookCertManager/Deployment: enabled with connectInject.enabled=true and controller.enabled=false" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'connectInject.enabled=true' \
      . | tee /dev/stderr |
      yq 'length > 0' | tee /dev/stderr)
  [ "${actual}" = "true" ]
}

@test "webhookCertManager/Deployment: enabled with connectInject.enabled=true and controller.enabled=true" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      . | tee /dev/stderr |
      yq 'length > 0' | tee /dev/stderr)
  [ "${actual}" = "true" ]
}

@test "webhookCertManager/Deployment: no tolerations by default" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      . | tee /dev/stderr |
      yq -r '.spec.template.spec.tolerations' | tee /dev/stderr)
  [ "${actual}" = "null" ]
}

@test "webhookCertManager/Deployment: tolerations can be set" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      --set 'webhookCertManager.tolerations=- key: value' \
      . | tee /dev/stderr |
      yq -r '.spec.template.spec.tolerations[0].key' | tee /dev/stderr)
  [ "${actual}" = "value" ]
}

@test "webhookCertManager/Deployment: no nodeSelector by default" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      . | tee /dev/stderr |
      yq -r '.spec.template.spec.nodeSelector' | tee /dev/stderr)
  [ "${actual}" = "null" ]
}

@test "webhookCertManager/Deployment: nodeSelector can be set" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      --set 'webhookCertManager.nodeSelector=- key: value' \
      . | tee /dev/stderr |
      yq -r '.spec.template.spec.nodeSelector[0].key' | tee /dev/stderr)
  [ "${actual}" = "value" ]
}

#--------------------------------------------------------------------
# Vault

@test "webhookCertManager/Deployment: disabled when the following are configured - global.secretsBackend.vault.enabled, global.secretsBackend.vault.enabled, global.secretsBackend.vault.connectInjectRole, global.secretsBackend.vault.connectInject.tlsCert.secretName, global.secretsBackend.vault.connectInject.caCert.secretName, global.secretsBackend.vault.controllerRole, global.secretsBackend.vault.controller.tlsCert.secretName, and .global.secretsBackend.vault.controller.caCert.secretName" {
  cd `chart_dir`
  assert_empty helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'global.secretsBackend.vault.enabled=true' \
      --set 'global.secretsBackend.vault.consulClientRole=test' \
      --set 'global.secretsBackend.vault.consulServerRole=foo' \
      --set 'global.secretsBackend.vault.consulCARole=carole' \
      --set 'global.secretsBackend.vault.connectInjectRole=inject-ca-role' \
      --set 'global.secretsBackend.vault.connectInject.tlsCert.secretName=pki/issue/connect-webhook-cert-dc1' \
      --set 'global.secretsBackend.vault.connectInject.caCert.secretName=pki/issue/connect-webhook-cert-dc1' \
      --set 'global.secretsBackend.vault.controllerRole=test' \
      --set 'global.secretsBackend.vault.controller.caCert.secretName=foo/ca' \
      --set 'global.secretsBackend.vault.controller.tlsCert.secretName=foo/tls' \
      --set 'global.secretsBackend.vault.consulClientRole=foo' \
      --set 'global.secretsBackend.vault.consulServerRole=bar' \
      --set 'global.secretsBackend.vault.consulCARole=test2' \
      .
}

#--------------------------------------------------------------------
# global.podSecurityStandards

@test "webhookCertManager/Deployment: podSecurityStandards default off" {
  cd `chart_dir`
  local actual=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      . | tee /dev/stderr |
      yq -r '.spec.template.spec.containers | map(select(.name == "webhook-cert-manager")) | .[0].securityContext | length > 0' | tee /dev/stderr)
  [ "${actual}" = "false" ]
}

@test "webhookCertManager/Deployment: global.podSecurityStandards are not set when global.openshift.enabled=true" {
  cd `chart_dir`
  local manifest=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      --set 'global.podSecurityStandards.securityContext.bob=false' \
      --set 'global.podSecurityStandards.securityContext.alice=true' \
      --set 'global.openshift.enabled=true' \
      . | tee /dev/stderr)

  local actual=$(echo "$manifest" | yq -r '.spec.template.spec.containers | map(select(.name == "webhook-cert-manager")) | .[0].securityContext')
  [ "${actual}" = "null" ]
}

@test "webhookCertManager/Deployment: global.podSecurityStandards can be set with tls and acls enabled" {
  cd `chart_dir`
  local object=$(helm template \
      -s templates/webhook-cert-manager-deployment.yaml  \
      --set 'controller.enabled=true' \
      --set 'connectInject.enabled=true' \
      --set 'global.podSecurityStandards.securityContext.bob=false' \
      --set 'global.podSecurityStandards.securityContext.alice=true' \
      . | tee /dev/stderr |
      yq '.spec.template.spec' | tee /dev/stderr)

  local actual=$(echo $object |
      yq -r '.containers | map(select(.name=="webhook-cert-manager")) | .[0].securityContext' | jq -r .bob)
  [ "${actual}" = "false" ]
  local actual=$(echo $object |
      yq -r '.containers | map(select(.name=="webhook-cert-manager")) | .[0].securityContext' | jq -r .alice)
  [ "${actual}" = "true" ]
}
