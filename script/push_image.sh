#!/bin/bash

REPO=swr.cn-south-1.myhuaweicloud.com/dc2gz

IMAGE=ops-server
TAG=$1

if [ -z "$TAG" ]; then
    TAG="dev"
fi

docker tag ${IMAGE}:${TAG} ${repo}/${IMAGE}:${TAG}
docker push ${repo}/${IMAGE}:${TAG}