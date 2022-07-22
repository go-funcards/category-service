#!/bin/sh

GOPATH=$(go env GOPATH)
PATH=$PATH:$GOPATH/bin

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protoc -I proto --go_out=proto/v1 --go-grpc_out=proto/v1 proto/v1/category.proto
