#!/bin/sh

# default values
origin_tag="local"

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

if [ $# -gt 0 ]; then
  origin_tag="$1"
fi

### main

# Executor nodes run inside containers using docker.
# They share access to the docker daemon on their host
# so we can use docker inside containers, but we can't
# use podman. Sharing the images is done through a
# local registry.

runcheck "failure when building redis image"  \
  docker build                                \
  -t shorturl/redis:${origin_tag}             \
  -f src/redis/Dockerfile                     \
  src/redis/

runcheck "failure when building worker image" \
  docker build                                \
  -t shorturl/worker:${origin_tag}            \
  -f src/worker/Dockerfile                    \
  src/worker/
  
#runcheck "failure when building server image" \
#  docker build                                \
#  -t shorturl/server:${origin_tag}            \
#  -f src/server/Dockerfile                    \
#  src/server/
