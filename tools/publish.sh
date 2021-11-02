#!/bin/sh

# default values
origin_tag="local"
target_tag="published"
registry=""

# default arch value is based on machine name
arch="$(uname -m)"
case "$arch" in
  x86_64)
    arch="arm64"
    ;;
  aarch64)
    arch="arm64"
    ;;
  *)
    arch="noarch"
    ;;
esac

# macros

# runs a command and checks for errors
# usage:
# runcheck "error message" command with options
runcheck() {
  if [ $# -lt 2 ]; then
    echo "missing arguments in macro runcheck"
    exit 1
  fi
  error_message="$1"
  shift
  echo "run: '$@'"
  eval "$@"
  if [ $? -ne 0 ]; then
    echo "${error_message}"
    exit 1
  fi
}

### main

if [ $# -gt 3 ]; then
  origin_tag="$1"
  target_tag="$2"
  arch="$3"
  registry="$4/"
fi

# Executor nodes run inside containers using docker.
# They share access to the docker daemon on their host
# so we can use docker inside containers, but we can't
# use podman. Sharing the images is done through a
# local registry.

runcheck "error while tagging redis image"  \
  docker tag                                \
  shorturl/redis:${origin_tag}              \
  ${registry}shorturl/redis:${target_tag}-${arch}

runcheck "error while pushing redis image"  \
  docker push                               \
  ${registry}shorturl/redis:${target_tag}-${arch}

runcheck "error while tagging worker image" \
  docker tag                                \
  shorturl/worker:${origin_tag}             \
  ${registry}shorturl/worker:${target_tag}-${arch}

runcheck "error while pushing worker image" \
  docker push                               \
  ${registry}shorturl/worker:${target_tag}-${arch}

#runcheck "error while tagging server image" \
#  docker tag                                \
#  shorturl/server:${origin_tag}             \
#  ${registry}shorturl/server:${target_tag}-${arch}

#runcheck "error while pushing server image" \
#  docker push                               \
#  ${registry}shorturl/server:${target_tag}-${arch}
