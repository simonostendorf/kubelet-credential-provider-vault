#!/bin/bash

set -e
set -o pipefail

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

REQUEST=$(jq ".serviceAccountToken=\""$(cat ./service-account-token.tmp)"\"" ./request.json)

RESPONSE=$(echo $REQUEST | ../../bin/kubelet-credential-provider-vault_${GOOS}_${GOARCH} \
  --log-level=debug \
  --vault-addr=${VAULT_ADDR:-"http://127.0.0.1:8200"} \
  --vault-auth-method=kubernetes \
  --vault-auth-mount=kubernetes \
  --vault-auth-role=example \
  --vault-secret-mount=secret \
  --vault-secret-path=example)

GOT_RESPONSE=$(jq -r '.' <<< $RESPONSE)
EXPECTED_RESPONSE=$(jq -r '.' response.json)

if [ "$GOT_RESPONSE" == "$EXPECTED_RESPONSE" ]; then
  echo "Test passed!"
else
  echo "Test failed!"
  echo "Expected: $EXPECTED_RESPONSE"
  echo "Got: $GOT_RESPONSE"
  exit 1
fi
