#!/bin/sh
set -e
#  clean proto
rm -rf ./pkg/grpc/pb/*

# generate proto
protoc -I ./pkg/grpc/proto/  ./pkg/grpc/proto/*.proto   --go-grpc_out=../ --go_out=../
protoc -I . ./internal/apps/vpn/proto/*.proto   --go-grpc_out=../ --go_out=../
protoc -I . ./internal/apps/auth/proto/*.proto   --go-grpc_out=../ --go_out=../
protoc-go-inject-tag -input ./pkg/grpc/pb/*/*.go
