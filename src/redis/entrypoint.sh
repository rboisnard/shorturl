#!/bin/sh

# default values
if [ -z ${REDIS_PORT} ]; then REDIS_PORT=6379; fi

if [ "x${BIND_IP}" = "x1" ]; then
  IP_OPTION="--bind $(hostname -i)"
fi
redis-server --port ${REDIS_PORT} ${IP_OPTION}

# TODO: demonize properly
while [ 1 -eq 1 ]; do
  sleep 30
done
