#!/usr/bin/env bash
IMAGE=pht_seneor
VERSION=0.618
docker rmi $IMAGE:$VERSION
docker build -t $IMAGE:$VERSION .
