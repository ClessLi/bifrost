#!/usr/bin/env bash

# shellcheck disable=SC2034 # Variables sourced in other scripts.

# The server platform we are building on.
readonly BIFROST_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
  windows/amd64
)

# The set of server targets that we are only building for Linux
# If you update this list, please also update build/BUILD.
bifrost::golang::server_targets() {
  local targets=(
    bifrost
  )
  echo "${targets[@]}"
}

IFS=" " read -ra BIFROST_SERVER_TARGETS <<< "$(bifrost::golang::server_targets)"
readonly BIFROST_SERVER_TARGETS
readonly BIFROST_SERVER_BINARIES=("${BIFROST_SERVER_TARGETS[@]##*/}")

# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use bifrost::util::read-array:
# bifrost::util::read-array FOO < <(bifrost::golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# Returns a sorted newline-separated list containing only duplicated items.
bifrost::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
bifrost::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Depends on values of user-facing BIFROST_BUILD_PLATFORMS, BIFROST_FASTBUILD,
# and BIFROST_BUILDER_OS.
# Configures BIFROST_SERVER_PLATFORMS and BIFROST_CLIENT_PLATFORMS, then sets them
# to readonly.
# The configured vars will only contain platforms allowed by the
# BIFROST_SUPPORTED* vars at the top of this file.
declare -a BIFROST_SERVER_PLATFORMS
bifrost::golang::setup_platforms() {
  if [[ -n "${BIFROST_BUILD_PLATFORMS:-}" ]]; then
    # BIFROST_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    local -a platforms
    IFS=" " read -ra platforms <<< "${BIFROST_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with bifrost::golang::dups
    # is not defeated by duplicates in user input.
    bifrost::util::read-array platforms < <(bifrost::golang::dedup "${platforms[@]}")

    # Use bifrost::golang::dups to restrict the builds to the platforms in
    # BIFROST_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    bifrost::util::read-array BIFROST_SERVER_PLATFORMS < <(bifrost::golang::dups \
        "${platforms[@]}" \
        "${BIFROST_SUPPORTED_SERVER_PLATFORMS[@]}" \
      )
    readonly BIFROST_SERVER_PLATFORMS

  elif [[ "${BIFROST_FASTBUILD:-}" == "true" ]]; then
    BIFROST_SERVER_PLATFORMS=(linux/amd64)
    readonly BIFROST_SERVER_PLATFORMS
  else
    BIFROST_SERVER_PLATFORMS=("${BIFROST_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly BIFROST_SERVER_PLATFORMS

  fi
}

bifrost::golang::setup_platforms

readonly BIFROST_ALL_TARGETS=(
  "${BIFROST_SERVER_TARGETS[@]}"
)
readonly BIFROST_ALL_BINARIES=("${BIFROST_ALL_TARGETS[@]##*/}")

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
bifrost::golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

# Ensure the go tool exists and is a viable version.
bifrost::golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    bifrost::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.13.4
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    bifrost::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
BIFROST requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# bifrost::golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# BIFROST build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
bifrost::golang::setup_env() {
  bifrost::golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}
