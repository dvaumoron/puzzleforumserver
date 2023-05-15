#!/usr/bin/env bash

./build/build.sh

buildah from --name puzzleforumserver-working-container scratch
buildah copy puzzleforumserver-working-container $HOME/go/bin/puzzleforumserver /bin/puzzleforumserver
buildah config --env SERVICE_PORT=50051 puzzleforumserver-working-container
buildah config --port 50051 puzzleforumserver-working-container
buildah config --entrypoint '["/bin/puzzleforumserver"]' puzzleforumserver-working-container
buildah commit puzzleforumserver-working-container puzzleforumserver
buildah rm puzzleforumserver-working-container

buildah push puzzleforumserver docker-daemon:puzzleforumserver:latest
