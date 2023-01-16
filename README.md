# svc-todo-gRPC

## Running database MYSQL with setting in /server/config/config.json and redis with setting in /server/db/rdb.go:

## Install protocol buffer for go
**Using this command':**
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
**Update PATH so that the protoc compiler can find the plugins::**
```
$ export PATH="$PATH:$(go env GOPATH)/bin"
```
**Create pb folder**
**Generate gRPC code using MAKEFILE**
```
$ make -f MakeFile gen-proto
```
**Folder proto will contain grpc.pb.go and pb.go**
**Run this command to update go module:**
```
$ go mod tidy
```
**Running golang in folder /server/main.go:**
**Run the following command:**
```
$ go run main.go
```
**Testing grpc request using postman in 0.0.0.0:50051**