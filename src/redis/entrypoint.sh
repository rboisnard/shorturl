#!/bin/sh

# default values
if [ -z ${REDIS_PORT} ]; then REDIS_PORT=6379; fi

redis-server --port ${REDIS_PORT}

while [ 1 -eq 1 ]; then
  sleep 30
done
