#!/bin/sh

# This script builds and pushes a Docker image to a registry.
# It is designed to be run as a command in CircleCI, a continuous integration and delivery platform.

# Enable the 'exit on error' option
set -e

# Set variables
# $1 is the first command line argument passed to the script (the image name)
# $2 is the second command line argument passed to the script (the root directory)
image=$1
root=$2

# Log in to Docker using the DOCKER_ID and DOCKER_PASSWD environment variables
echo "$DOCKER_PASSWD" | docker login -u "$DOCKER_ID" --password-stdin

# Change to the root directory
cd "$root"

# Build and push Docker image
# --build-arg sets a build-time variable for the Docker image
# -t specifies the name and tag for the image
docker build --build-arg VERSION="$CIRCLE_SHA1" --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" -t "$image:$CIRCLE_SHA1" .
docker push "$image:$CIRCLE_SHA1"

# If a tag is present, build and push an additional image with the tag
if [ -n "$CIRCLE_TAG" ]; then
  docker build --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" -t "$image:$CIRCLE_TAG" .
  docker push "$image:$CIRCLE_TAG"
fi

# Log out of Docker
docker logout
