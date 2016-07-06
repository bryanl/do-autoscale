#!/bin/bash

set -ev

PATH=$(pwd)/bin:$(pwd)/vendor/bin:$PATH

docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

cd vendor
gb build .../go-bindata
cd ..

if [ "$TRAVIS_BRANCH" = "master" ]; then
  make build-app
else
  DOCKER_TAG=${TRAVIS_BRANCH} BUILD_LATEST=0 make build-app
fi

docker push bryanl/do-autoscale