#!/bin/sh

set -e
set -u

SERVICE_NAME="${PWD##*/}"

CONTAINER_NAME="jmbarzee/services_$SERVICE_NAME" 

GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse HEAD)}
GIT_BRANCH=${GIT_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}



docker image build \
    --tag="$CONTAINER_NAME:$GIT_COMMIT" \
    .

docker image push "$CONTAINER_NAME:$GIT_COMMIT"

if  [[ $GIT_BRANCH == "main" ]]; then 
    docker image tag "$CONTAINER_NAME:$GIT_COMMIT" "$CONTAINER_NAME:latest"
    docker image push "$CONTAINER_NAME:latest"
fi 