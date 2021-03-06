#!/bin/sh
set -e

export dockerfile="Dockerfile"
export arch=$(uname -m)

export eTAG="latest"

if [ "$arch" = "armv7l" ] ; then
   dockerfile="Dockerfile.armhf"
   eTAG="latest-armhf-dev"
fi

echo "$1"
if [ "$1" ] ; then
  eTAG=$1
  if [ "$arch" = "armv7l" ] ; then
    eTAG="$1-armhf"
  fi
fi

NS=lambdanic
REPO=faas-gateway

echo Building $NS/$REPO:$eTAG

GIT_COMMIT_MESSAGE=$(git log -1 --pretty=%B 2>&1 | head -n 1)
GIT_COMMIT_SHA=$(git rev-list -1 HEAD)
VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///' || echo dev)

while true; do
    read -p "This will delete all docker images in your system. Continue? " yn
    case $yn in
        [Yy]* ) docker rmi $(docker images -a -q) || true; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes[Y] or no[N].";;
    esac
done

dep ensure -update

docker build --build-arg https_proxy=$https_proxy --build-arg http_proxy=$http_proxy \
  --build-arg GIT_COMMIT_MESSAGE="$GIT_COMMIT_MESSAGE" --build-arg GIT_COMMIT_SHA=$GIT_COMMIT_SHA \
  --build-arg VERSION=${VERSION:-dev} \
  -t $NS/$REPO:$eTAG . -f $dockerfile --no-cache
