#!/bin/bash

set -e
set -o pipefail

# create vault and k3s
docker compose up -d

# wait for vault to be ready
while ! curl -s ${VAULT_ADDR:-"http://127.0.0.1:8200"}/v1/sys/health | grep -q '"sealed":false'; do
  echo "Waiting for Vault to be unsealed..."
  sleep 5
done
echo "Vault is ready and unsealed."

# wait for k3s to be ready
while ! curl -sk ${KUBE_HOST:-"https://127.0.0.1:6443"}/readyz | grep -q 'Unauthorized'; do
  echo "Waiting for K3s to be ready..."
  sleep 5
done
echo "K3s is ready."

# setup vault and k3s
tofu -chdir=./setup init
tofu -chdir=./setup apply -auto-approve
TOFU_OUTPUTS=$(tofu -chdir=./setup output -json)

# read service account token and save to file
jq -r '.service_account_token.value' <<< "$TOFU_OUTPUTS" > ./service-account-token.tmp
