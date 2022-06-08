#!/usr/bin/env bash
IMAGE=pht_sensor
VERSION=0.9
docker rmi $IMAGE:$VERSION
docker build -t $IMAGE:$VERSION .
