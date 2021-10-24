#!/bin/sh

podman run --rm -d -p 5500:5500 -e PORT=5500 shorturl/worker:staging