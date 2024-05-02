# gRPC echo server

This is a simple ping pong service built using Go and gRPC.

![Application](./ping-pong-app.png)

## Build Images and deploy the application

- `go build -ldflags="-s -w" -o ping-pong ./cmd/...`
- `sudo ./scripts/build_images.sh -u username -t latest`
- `kubectl apply -f ping-pong-app.yaml`


## Send queries

- `curl http://10.96.88.88:8080/ping-echo?body=hello`
- `./wrk/wrk -d 10s -c 1 -t 1 http://10.96.88.88:8080/ping-echo -s ./application/wrk_scripts/echo.lua -L`