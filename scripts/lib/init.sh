#!/usr/bin/env bash

set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
BIFROST_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${BIFROST_ROOT}/scripts/lib/util.sh"
source "${BIFROST_ROOT}/scripts/lib/logging.sh"
source "${BIFROST_ROOT}/scripts/lib/color.sh"

bifrost::log::install_errexit

source "${BIFROST_ROOT}/scripts/lib/version.sh"
source "${BIFROST_ROOT}/scripts/lib/golang.sh"