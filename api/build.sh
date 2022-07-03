#!/bin/bash

# Ensure you have grpc and protobuf installed before 
# executing this script
# brew install protobuf
# brew install grpc
# brew install protoc-gen-go
# brew install protoc-gen-go-grpc

protoc api/v1/service.proto \
--go_out=. \
--go_opt=paths=source_relative \
--proto_path=.

protoc api/v1/service.proto \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--proto_path=.