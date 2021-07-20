#!/bin/bash

.PHONY:
gen: clean
	@ echo "generating protobuf"
	@ protoc --proto_path=proto proto/*.proto \
 	  --go_out=gen/pb --go-grpc_out=gen/pb \
 	  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

.PHONY:
fmt:
	@ echo "formatting go code"
	@ go fmt ./...

.PHONY:
server:
	@ echo "running server"
	@ go run cmd/server/main.go -port 8080

.PHONY:
client:
	@ echo "running client"
	@ go run cmd/client/main.go -address 0.0.0.0:8080

.PHONY:
clean:
	@ echo "cleaning-up"
	@ rm -rf gen/pb/*.pb.go