#!/bin/sh

podman build -t shorturl/redis:local -f src/redis/Dockerfile src/redis/
podman build -t shorturl/worker:local -f src/worker/Dockerfile src/worker/