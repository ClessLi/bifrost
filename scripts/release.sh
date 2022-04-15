#!/usr/bin/env bash

# Build a Bifrost release.  This will build the binaries, create the Docker
# images and other build artifacts.

set -o errexit
set -o nounset
set -o pipefail

#set -x
BIFROST_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${BIFROST_ROOT}/scripts/common.sh"
source "${BIFROST_ROOT}/scripts/lib/release.sh"

BIFROST_RELEASE_RUN_TESTS=${BIFROST_RELEASE_RUN_TESTS-y}

bifrost::golang::setup_env
bifrost::build::verify_prereqs
bifrost::release::verify_prereqs
#bifrost::build::build_image
bifrost::build::build_command
bifrost::release::package_tarballs
bifrost::release::github_release
bifrost::release::generate_changelog