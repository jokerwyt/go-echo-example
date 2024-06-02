#!/bin/bash
set -ex

# # Ensure the script is being run from the correct directory
# SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
# cd "$SCRIPT_DIR"

go build -o server server.go
go build -o frontend frontend.go


# Build the frontend image
sudo docker build --tag echo-frontend -f Dockerfile-frontend ../..

# Build the server image
sudo docker build --tag echo-server -f Dockerfile-server ../..

# Tag the images
sudo docker tag echo-frontend  <docker-username>/echo-frontend-grpc:namespaced
sudo docker tag echo-server  <docker-username>/echo-server-grpc:namespaced

# Push the images to the registry
sudo docker push  <docker-username>/echo-frontend-grpc:namespaced
sudo docker push  <docker-username>/echo-server-grpc:namespaced

set +ex
