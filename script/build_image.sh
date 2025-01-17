#!/bin/bash

IMAGE=ops-server
TAG=$1

if [ -z "$TAG" ]; then
    TAG="dev"
fi

docker build -t ${IMAGE}:${TAG} .