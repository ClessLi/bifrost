#!/usr/bin/env bash

# This file creates release artifacts (tar files, container images) that are
# ready to distribute to install or distribute to end users.

###############################################################################
# Most of the ::release:: namespace functions have been moved to
# github.com/bifrost/release.  Have a look in that repo and specifically in
# lib/releaselib.sh for ::release::-related functionality.
###############################################################################

# This is where the final release artifacts are created locally
readonly RELEASE_STAGE="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_TARS="${LOCAL_OUTPUT_ROOT}/release-tars"

# Bifrost git branch info
readonly BIFROST_CURRENT_BRANCH="$(git branch --show-current)"

# BIFROST github account info
readonly BIFROST_GITHUB_ORG=ClessLi
readonly BIFROST_GITHUB_REPO=bifrost

readonly ARTIFACT=bifrost.tar.gz
readonly SERVER_ARTIFACTS='bifrost-server-*-*.tar.gz'
readonly CHECKSUM=${ARTIFACT}.sha1sum

BIFROST_BUILD_CONFORMANCE=${BIFROST_BUILD_CONFORMANCE:-y}

# Validate a ci version
#
# Globals:
#   None
# Arguments:
#   version
# Returns:
#   If version is a valid ci version
# Sets:                    (e.g. for '1.2.3-alpha.4.56+abcdef12345678')
#   VERSION_MAJOR          (e.g. '1')
#   VERSION_MINOR          (e.g. '2')
#   VERSION_PATCH          (e.g. '3')
#   VERSION_PRERELEASE     (e.g. 'alpha')
#   VERSION_PRERELEASE_REV (e.g. '4')
#   VERSION_BUILD_INFO     (e.g. '.56+abcdef12345678')
#   VERSION_COMMITS        (e.g. '56')
function bifrost::release::parse_and_validate_ci_version() {
  # Accept things like "v1.2.3-alpha.4.56+abcdef12345678" or "v1.2.3-beta.4"
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-zA-Z0-9]+)\\.(0|[1-9][0-9]*)(\\.(0|[1-9][0-9]*)\\+[0-9a-f]{7,40})?$"
  local -r version="${1-}"
  [[ "${version}" =~ ${version_regex} ]] || {
    bifrost::log::error "Invalid ci version: '${version}', must match regex ${version_regex}"
    return 1
  }

  # The VERSION variables are used when this file is sourced, hence
  # the shellcheck SC2034 'appears unused' warning is to be ignored.

  # shellcheck disable=SC2034
  VERSION_MAJOR="${BASH_REMATCH[1]}"
  # shellcheck disable=SC2034
  VERSION_MINOR="${BASH_REMATCH[2]}"
  # shellcheck disable=SC2034
  VERSION_PATCH="${BASH_REMATCH[3]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE="${BASH_REMATCH[4]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE_REV="${BASH_REMATCH[5]}"
  # shellcheck disable=SC2034
  VERSION_BUILD_INFO="${BASH_REMATCH[6]}"
  # shellcheck disable=SC2034
  VERSION_COMMITS="${BASH_REMATCH[7]}"
}

# ---------------------------------------------------------------------------
# Build final release artifacts
function bifrost::release::clean_cruft() {
  # Clean out cruft
  find "${RELEASE_STAGE}" -name '*~' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '#*#' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '.DS*' -exec rm {} \;
}

function bifrost::release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
#  bifrost::release::package_src_tarball &
  bifrost::release::package_server_tarballs &
  bifrost::util::wait-for-jobs || { bifrost::log::error "previous tarball phase failed"; return 1; }

  bifrost::release::package_final_tarball & # _final depends on some of the previous phases
  bifrost::util::wait-for-jobs || { bifrost::log::error "previous tarball phase failed"; return 1; }
}

# Package the source code we built, for compliance/licensing/audit/yadda.
function bifrost::release::package_src_tarball() {
  local -r src_tarball="${RELEASE_TARS}/bifrost-src.tar.gz"
  bifrost::log::status "Building tarball: src"
  if [[ "${BIFROST_GIT_TREE_STATE-}" = 'clean' ]]; then
    git archive -o "${src_tarball}" HEAD
  else
    find "${BIFROST_ROOT}" -mindepth 1 -maxdepth 1 \
      ! \( \
      \( -path "${BIFROST_ROOT}"/_\* -o \
      -path "${BIFROST_ROOT}"/.git\* -o \
      -path "${BIFROST_ROOT}"/.gitignore\* -o \
      -path "${BIFROST_ROOT}"/.gsemver.yaml\* -o \
      -path "${BIFROST_ROOT}"/.config\* -o \
      -path "${BIFROST_ROOT}"/.chglog\* -o \
      -path "${BIFROST_ROOT}"/.gitlint -o \
      -path "${BIFROST_ROOT}"/.golangci.yaml -o \
      -path "${BIFROST_ROOT}"/.goreleaser.yml -o \
      -path "${BIFROST_ROOT}"/.note.md -o \
      -path "${BIFROST_ROOT}"/.todo.md \
      \) -prune \
      \) -print0 \
      | "${TAR}" czf "${src_tarball}" --transform "s|${BIFROST_ROOT#/*}|bifrost|" --null -T -
  fi
}

# Package up all of the server binaries
function bifrost::release::package_server_tarballs() {
  # Find all of the built client binaries
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  if [[ -n ${BIFROST_BUILD_PLATFORMS-} ]]; then
    read -ra long_platforms <<< "${BIFROST_BUILD_PLATFORMS}"
  fi

  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    bifrost::log::status "Starting tarball: server $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/server/${platform_tag}/bifrost"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    local server_bins=("${BIFROST_SERVER_BINARIES[@]}")

      # This fancy expression will expand to prepend a path
      # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
      # server_bins array.
      cp "${server_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
        "${release_stage}/server/bin/"

      bifrost::release::clean_cruft

      local package_name="${RELEASE_TARS}/bifrost-server-${platform_tag}.tar.gz"
      bifrost::release::create_tarball "${package_name}" "${release_stage}/.."
      ) &
    done

    bifrost::log::status "Waiting on tarballs"
    bifrost::util::wait-for-jobs || { bifrost::log::error "server tarball creation failed"; exit 1; }
}


function bifrost::release::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function bifrost::release::sha1() {
  if which sha1sum >/dev/null 2>&1; then
    sha1sum "$1" | awk '{ print $1 }'
  else
    shasum -a1 "$1" | awk '{ print $1 }'
  fi
}

# This is all the platform-independent stuff you need to run/install bifrost.
# Arch-specific binaries will need to be downloaded separately (possibly by
# using the bundled cluster/get-bifrost-binaries.sh script).
# Included in this tarball:
#   - Cluster spin up/down scripts and configs for various cloud providers
#   - Tarballs for manifest configs that are ready to be uploaded
#   - Examples (which may or may not still work)
#   - The remnants of the docs/ directory
function bifrost::release::package_final_tarball() {
  bifrost::log::status "Building tarball: final"

  # This isn't a "full" tarball anymore, but the release lib still expects
  # artifacts under "full/bifrost/"
  local release_stage="${RELEASE_STAGE}/full/bifrost"
  rm -rf "${release_stage}"
  mkdir -p "${release_stage}"

  # We want everything in /scripts.
  mkdir -p "${release_stage}/release"
  cp -R "${BIFROST_ROOT}/scripts/release" "${release_stage}/"

  mkdir -p "${release_stage}/server"
#  cat <<EOF > "${release_stage}/server/README"
#Server binary tarballs are no longer included in the Bifrost final tarball.
#EOF
  cp ${RELEASE_TARS}/${SERVER_ARTIFACTS} ${release_stage}/server

  # Include hack/lib as a dependency for the cluster/ scripts
  #mkdir -p "${release_stage}/hack"
  #cp -R "${BIFROST_ROOT}/hack/lib" "${release_stage}/hack/"

  cp -R ${BIFROST_ROOT}/{docs,configs,scripts,init,README.md,LICENSE} "${release_stage}/"

  echo "${BIFROST_GIT_VERSION}" > "${release_stage}/version"

  bifrost::release::clean_cruft

  local package_name="${RELEASE_TARS}/${ARTIFACT}"
  bifrost::release::create_tarball "${package_name}" "${release_stage}/.."
}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/bifrostis what is
# being packaged.
function bifrost::release::create_tarball() {
  bifrost::build::ensure_tar

  local tarfile=$1
  local stagingdir=$2

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" bifrost --owner=0 --group=0
}

function bifrost::release::install_github_release(){
  GO111MODULE=on go install github.com/github-release/github-release@latest
}

# Require the following tools:
# - github-release
# - gsemver
# - git-chglog
# - coscmd
function bifrost::release::verify_prereqs(){
  if [ -z "$(which github-release 2>/dev/null)" ]; then
    bifrost::log::info "'github-release' tool not installed, try to install it."

    if ! bifrost::release::install_github_release; then
      bifrost::log::error "failed to install 'github-release'"
      return 1
    fi
  fi

  if [ -z "$(which git-chglog 2>/dev/null)" ]; then
    bifrost::log::info "'git-chglog' tool not installed, try to install it."

    if ! go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest &>/dev/null; then
      bifrost::log::error "failed to install 'git-chglog'"
      return 1
    fi
  fi

  if [ -z "$(which gsemver 2>/dev/null)" ]; then
    bifrost::log::info "'gsemver' tool not installed, try to install it."

    if ! go install github.com/arnaud-deprez/gsemver@latest &>/dev/null; then
      bifrost::log::error "failed to install 'gsemver'"
      return 1
    fi
  fi

}

# Create a github release with specified tarballs.
# NOTICE: Must export 'GITHUB_TOKEN' env in the shell, details:
# https://github.com/github-release/github-release
function bifrost::release::github_release() {
  # create a github release
  set -x
  bifrost::log::info "create a new github release with tag ${BIFROST_GIT_VERSION}"
  github-release release \
    --pre-release \
    --user ${BIFROST_GITHUB_ORG} \
    --repo ${BIFROST_GITHUB_REPO} \
    --tag ${BIFROST_GIT_VERSION} \
    --description "${BIFROST_CHANGELOG}"

  set +x
  # update bifrost tarballs
  bifrost::log::info "upload ${ARTIFACT} to release ${BIFROST_GIT_VERSION}"
  github-release upload \
    --user ${BIFROST_GITHUB_ORG} \
    --repo ${BIFROST_GITHUB_REPO} \
    --tag ${BIFROST_GIT_VERSION} \
    --name ${ARTIFACT} \
    --file ${RELEASE_TARS}/${ARTIFACT}

#  bifrost::log::info "upload bifrost-src.tar.gz to release ${BIFROST_GIT_VERSION}"
#  github-release upload \
#    --user ${BIFROST_GITHUB_ORG} \
#    --repo ${BIFROST_GITHUB_REPO} \
#    --tag ${BIFROST_GIT_VERSION} \
#    --name "bifrost-src.tar.gz" \
#    --file ${RELEASE_TARS}/bifrost-src.tar.gz
}

function bifrost::release::generate_changelog() {
  bifrost::log::info "generate CHANGELOG-${BIFROST_GIT_VERSION#v}.md and commit it"

  local CHANGELOG_COMMIT_MSG="docs(changelog): add \`CHANGELOG-${BIFROST_GIT_VERSION#v}.md\`"

  if [[ "$(git log origin/${BIFROST_CURRENT_BRANCH} | grep -F "${CHANGELOG_COMMIT_MSG}" | wc -l)" -ne 0 ]]
    then
      bifrost::log::info "CHANGELOG-${BIFROST_GIT_VERSION#v}.md has been committed and pushed"
      return
  fi

  echo "${BIFROST_CHANGELOG}" > ${BIFROST_ROOT}/CHANGELOG/CHANGELOG-${BIFROST_GIT_VERSION#v}.md

#  current_commit_id=$(git log HEAD -n 1 --pretty=format:%H)
  git add ${ROOT_DIR}/CHANGELOG/CHANGELOG-${BIFROST_GIT_VERSION#v}.md
  git commit -a -m "${CHANGELOG_COMMIT_MSG}"
  if ! git push origin HEAD
    then
#      git reset --soft "${current_commit_id}"
      bifrost::log::error_exit "failed to push commit"
  fi
}

