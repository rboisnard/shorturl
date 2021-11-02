#!/bin/sh

# default values
target_tag="published"
archs="amd64,arm64"
registry=""

# macros

# runs a command and checks for errors
# usage:
# runcheck "error message" command with options
runcheck() {
  if [ $# -lt 2 ]; then
    echo "missing arguments in macro runcheck"
    echo "(in: runcheck $@)"
    exit 1
  fi
  error_message="$1"
  shift
  echo "++run: $@"
  eval "$@"
  if [ $? -ne 0 ]; then
    echo "${error_message}"
    exit 1
  fi
}

# prints a command and runs it
# usage:
# runshow command with options
runshow() {
  echo "++run: $@"
  eval "$@"
}

# creates and publishes a manifest for multiple archs
# usage:
# create_manifest <registry/image:tag> <arch1>[,arch2][,arch3]
create_manifest() {
  if [ $# -lt 2 ]; then
    echo "missing arguments in macro create_manifest"
    echo "(in: create_manifest $@)"
    exit 1
  fi

  registry_image_tag="$1"
  archs="$2"

  DOCKER_CLI_EXPERIMENTAL=enabled

  amend_opts=""
  for arch in $(echo "${archs}" | tr "," " "); do
    # pull image for each arch
    runcheck "error when pulling ${registry_image_tag}-${arch}" \
      docker pull ${registry_image_tag}-${arch}

    # add arch to manifest option
    runshow "amend_opts=\"\${amend_opts} --amend ${registry_image_tag}-${arch}\""
  done

  # create manifest
  runcheck "error when creating manifest for ${registry_image_tag}" \
    docker manifest create --insecure ${registry_image_tag} ${amend_opts}

  # push manifest
  runcheck "error when pushing manifest for ${registry_image_tag}"  \
    docker manifest push --insecure ${registry_image_tag}

  # check manifest
  runcheck "error when checking manifest for ${registry_image_tag}" \
    docker manifest inspect --insecure ${registry_image_tag}
}

### main

if [ $# -gt 1 ]; then
  target_tag="$1"
  registry="$2/"
  archs="$3"
fi

create_manifest ${registry}shorturl/redis:${target_tag} ${archs}
create_manifest ${registry}shorturl/worker:${target_tag} ${archs}
#create_manifest ${registry}shorturl/server:${target_tag} ${archs}
