package grpc

// This file is exclusively used for hooking gRPC code generation into go
//go:generate protoc --proto_path=./proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative light_orchestrator.proto
