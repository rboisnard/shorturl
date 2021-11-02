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
    echo "(in: runcheck $@)"
    cleanup_containers bypass
    exit 1
  fi
  error_message="$1"
  shift
  echo "++run: $@"
  eval "$@"
  if [ $? -ne 0 ]; then
    echo "${error_message}"
    cleanup_containers bypass
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

# add arg "bypass" to avoid exiting with an error when the stop failed
cleanup_containers() {
  bypass=0
  if [ $# -gt 0 ] && [ "x$1" = "xbypass" ]; then
    bypass=1
  fi

  cleanup_error=0

  # stop worker container if it exists
  docker ps -a | grep shorturl_worker_${suffix}
  if [ $? -eq 0 ]; then
    runshow docker stop shorturl_worker_${suffix}
    if [ $? -ne 0 ]; then
      echo "error while stopping worker container"
      cleanup_error=1
    fi
  fi

  # stop redis container if it exists
  docker ps -a | grep shorturl_redis_${suffix}
  if [ $? -eq 0 ]; then
    runshow docker stop shorturl_redis_${suffix}
    if [ $? -ne 0 ]; then
      echo "error while stopping redis container"
      cleanup_error=1
    fi
  fi

  if [ ${bypass} -eq 0 ] && [ ${cleanup_error} -eq 1 ]; then
    exit 1
  fi
}

### main

if [ $# -gt 0 ]; then
  origin_tag="$1"
fi
suffix="${origin_tag}_$(date +%s)"

# start redis container
runcheck "cannot start redis container" \
  docker run --rm -d                    \
  -p 6379:6379                          \
  -e REDIS_PORT=6379                    \
  -e BIND_IP=1                          \
  --name=shorturl_redis_${suffix}       \
  shorturl/redis:${origin_tag}

# check redis container logs
runcheck "redis container failed during init" \
  "docker logs shorturl_redis_${suffix} | grep '# Server initialized'"

# get parameters for the worker
runshow "shorturl_redis_ip=\$(docker exec shorturl_redis_${suffix} hostname -i)"
runshow "host_ip=\$(ip -4 route show default | awk '{print \$3}')"

# start worker container
runcheck "cannot start worker container"  \
  docker run --rm -d                      \
  -p 8080:5500                            \
  -e APP_URL=${host_ip}:8080              \
  -e APP_PORT=5500                        \
  -e REDIS_IP=${shorturl_redis_ip}        \
  -e REDIS_PORT=6379                      \
  --name=shorturl_worker_${suffix}        \
  shorturl/worker:${origin_tag}

# check worker container logs
runcheck "worker container cannot access redis container" \
  "docker logs shorturl_worker_${suffix} | grep 'PONG <nil>'"

# TODO: test multiple requests
# test mock url (temporary)
runcheck "curl request is not getting the expected reply" \
  "curl -X POST ${host_ip}:8080 2> /dev/null | grep \"http://${host_ip}:8080/mock_url\""

echo "all tests ok"
cleanup_containers
