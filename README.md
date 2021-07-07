# grpcdemo-go

go语言简单实现gRPC示例。

## 快速开始
1. run rpc server `go run server/main.go`
2. run rpc client `go run client/main.go`

## 安装
```shell
go get google.golang.org/grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```