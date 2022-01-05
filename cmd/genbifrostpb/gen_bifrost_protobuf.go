package main

//go:generate protoc -I=../.. --go_out=plugins=grpc:../.. ../../api/protobuf-spec/bifrostpb/v1/bifrost.proto

func main() {
}
