#!/bin/bash
set -ex

sudo docker build --tag echo-frontend -f Dockerfile-frontend ..
sudo docker build --tag echo-server -f Dockerfile-server  ..
sudo docker tag echo-frontend nikolabo/echo-frontend-grpc
sudo docker tag echo-server nikolabo/echo-server-grpc
sudo docker push nikolabo/echo-frontend-grpc
sudo docker push nikolabo/echo-server-grpc

set +ex
