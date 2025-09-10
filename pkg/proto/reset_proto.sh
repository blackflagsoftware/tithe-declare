#! /bin/bash
cd ../.. && protoc --go_out=./pkg/proto --go-grpc_out=./pkg/proto ./pkg/proto/tithe-declare.proto