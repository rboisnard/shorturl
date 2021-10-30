#!/bin/sh

LOCAL_DIR=$(readlink -f $(dirname $(readlink -f $0))/..)

podman build -t raspi01:5000/shorturl/worker:dev -f src/worker/Dockerfile --target dev src/worker/
podman run --rm -it -v ${LOCAL_DIR}/src/worker/app:/shorturl -w /shorturl raspi01:5000/shorturl/worker:dev bash