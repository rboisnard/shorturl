#!/bin/sh

# default values
if [ -z ${APP_URL} ];     then APP_URL=localhost:8080;  fi
if [ -z ${APP_IP} ];      then APP_IP=localhost;        fi
if [ -z ${APP_PORT} ];    then APP_PORT=5500;           fi
if [ -z ${REDIS_IP} ];    then REDIS_IP=localhost;      fi
if [ -z ${REDIS_PORT} ];  then REDIS_PORT=6379;         fi

/app/shorturl
