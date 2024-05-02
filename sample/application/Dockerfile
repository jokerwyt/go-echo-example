# FROM golang:1.22.2-alpine as builder
# ENV CGO_ENABLED=1
# RUN apk add --no-cache --update go gcc g++
# # Set the working directory inside the container
# WORKDIR /app
# # Copy only the necessary Go mod files to cache dependencies
# COPY go.mod go.sum ./
# # Copy the AppNet library into the container
# COPY ./grpc/plugininterceptor /grpc/plugininterceptor
# # Copy the protobuf module into the container
# COPY ./ping_pb /ping_pb
# # Download and cache Go dependencies
# # RUN go mod download
# # Copy the entire project directory to the container
# COPY ./envoy .
# # Build the Go application with optimized flags
# RUN go build -ldflags="-s -w" -o /app/ping-pong ./cmd/...
# # Use a minimal base image for the final container
FROM alpine:latest
RUN apk add gcompat
# Set the working directory inside the container
WORKDIR /app
# Copy the built binary from the builder stage
COPY ./application/ping-pong ./ping-pong
# Set the entrypoint command to run the binary
CMD ["./ping-pong"]