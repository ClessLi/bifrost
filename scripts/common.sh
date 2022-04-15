#!/usr/bin/env bash

# shellcheck disable=SC2034 # Variables sourced in other scripts.

# Common utilities, variables and checks for all build scripts.
set -o errexit
set -o nounset
set -o pipefail

# Unset CDPATH, having it set messes up with script import paths
unset CDPATH

USER_ID=$(id -u)
GROUP_ID=$(id -g)

# This will canonicalize the path
BIFROST_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${BIFROST_ROOT}/scripts/lib/init.sh"

# Here we map the output directories across both the local and remote _output
# directories:
#
# *_OUTPUT_ROOT    - the base of all output in that environment.
# *_OUTPUT_SUBPATH - location where golang stuff is built/cached.  Also
#                    persisted across docker runs with a volume mount.
# *_OUTPUT_BINPATH - location where final binaries are placed.  If the remote
#                    is really remote, this is the stuff that has to be copied
#                    back.
# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${ROOT_DIR}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/platforms"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}"
readonly LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"

# ---------------------------------------------------------------------------
# Basic setup functions

# Verify that the right utilities and such are installed for building bifrost. Set
# up some dynamic constants.
# Args:
#   $1 - boolean of whether to require functioning docker (default true)
#
# Vars set:
#   BIFROST_ROOT_HASH

#   LOCAL_OUTPUT_BUILD_CONTEXT
function bifrost::build::verify_prereqs() {
  bifrost::log::status "Verifying Prerequisites...."
  bifrost::build::ensure_tar || return 1
  bifrost::build::ensure_rsync || return 1

  BIFROST_GIT_BRANCH=$(git symbolic-ref --short -q HEAD 2>/dev/null || true)
  BIFROST_ROOT_HASH=$(bifrost::build::short_hash "${HOSTNAME:-}:${BIFROST_ROOT}:${BIFROST_GIT_BRANCH}")

  bifrost::version::get_version_vars
  #bifrost::version::save_version_vars "${BIFROST_ROOT}/.dockerized-bifrost-version-defs"
  BIFROST_CHANGELOG="$(git-chglog ${BIFROST_GIT_VERSION})"
}

# ---------------------------------------------------------------------------
# Utility functions

function bifrost::build::is_gnu_sed() {
  [[ $(sed --version 2>&1) == *GNU* ]]
}

function bifrost::build::ensure_rsync() {
  if [[ -z "$(which rsync)" ]]; then
    bifrost::log::error "Can't find 'rsync' in PATH, please fix and retry."
    return 1
  fi
}

function  bifrost::build::set_proxy() {
  if [[ -n "${BIFROSTRNETES_HTTPS_PROXY:-}" ]]; then
    echo "ENV https_proxy $BIFROSTRNETES_HTTPS_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${BIFROSTRNETES_HTTP_PROXY:-}" ]]; then
    echo "ENV http_proxy $BIFROSTRNETES_HTTP_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${BIFROSTRNETES_NO_PROXY:-}" ]]; then
    echo "ENV no_proxy $BIFROSTRNETES_NO_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
}

function bifrost::build::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=tar
  if which gtar &>/dev/null; then
      TAR=gtar
  else
      if which gnutar &>/dev/null; then
	  TAR=gnutar
      fi
  fi
  if ! "${TAR}" --version | grep -q GNU; then
    echo "  !!! Cannot find GNU tar. Build on Linux or install GNU tar"
    echo "      on Mac OS X (brew install gnu-tar)."
    return 1
  fi
}

function bifrost::build::has_ip() {
  which ip &> /dev/null && ip -Version | grep 'iproute2' &> /dev/null
}

# Takes $1 and computes a short has for it. Useful for unique tag generation
function bifrost::build::short_hash() {
  [[ $# -eq 1 ]] || {
    bifrost::log::error "Internal error.  No data based to short_hash."
    exit 2
  }

  local short_hash
  if which md5 >/dev/null 2>&1; then
    short_hash=$(md5 -q -s "$1")
  else
    short_hash=$(echo -n "$1" | md5sum)
  fi
  echo "${short_hash:0:10}"
}

# ---------------------------------------------------------------------------
# Building


function bifrost::build::clean() {

  if [[ -d "${LOCAL_OUTPUT_ROOT}" ]]; then
    bifrost::log::status "Removing _output directory"
    rm -rf "${LOCAL_OUTPUT_ROOT}"
  fi
}

# Build all Bifrost commands.
function bifrost::build::build_command() {
  bifrost::log::status "Running build command..."
  make -C "${BIFROST_ROOT}" build.multiarch BINS="bifrost" PLATFORMS="linux_amd64 windows_amd64"
}
