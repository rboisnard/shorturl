#!/bin/sh

podman build -t raspi01:5000/shorturl/redis:staging -f src/redis/Dockerfile src/redis/
podman build -t raspi01:5000/shorturl/worker:staging -f src/worker/Dockerfile src/worker/

podman push raspi01:5000/shorturl/redis:staging
podman push raspi01:5000/shorturl/worker:staging