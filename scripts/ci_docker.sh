#!/bin/bash

set -ev

docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

if [ "$TRAVIS_BRANCH" = "master" ]; then
  make generate build-app
else
  DOCKER_TAG=$TRAVIS_BRANCH BUILD_LATEST=0 make generate build-app
fi

docker push bryanl/do-autoscale