#!/bin/bash

set -e
set -o pipefail

# destroy service account token file
rm ./service-account-token.tmp

# destroy log file
rm ./kubelet-credential-provider-vault.log

# destroy vault and k3s
docker compose down -v
