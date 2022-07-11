package main

import category "github.com/go-funcards/category-service/cmd"

//go:generate protoc -I proto --go_out=./proto/v1 --go-grpc_out=./proto/v1 proto/v1/category.proto

func main() {
	category.Execute()
}
