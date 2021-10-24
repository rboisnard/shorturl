#!/bin/sh

podman build -t shorturl/worker:staging -f src/worker/Dockerfile src/worker/