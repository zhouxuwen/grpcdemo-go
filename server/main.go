package main

import (
	"errors"
	"fmt"
	"grpcdemo-go/protobuf"
	"grpcdemo-go/server/interceptor"
	"grpcdemo-go/server/service"
	"net"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func main() {
	address := "localhost:9000"
	opts := make([]grpc.ServerOption, 0)

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
	if err != nil {
		fmt.Printf("failed to generate credentials err:%v\n", err)
		return
	}
	opts = append(opts, grpc.Creds(creds))

	// token auth
	opts = append(opts, grpc.UnaryInterceptor(interceptor.Interceptor))

	// rpc server
	server := grpc.NewServer(opts...)
	protobuf.RegisterHelloServiceServer(server, &service.HelloService)

	// tcp listen
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("failed to listen: %v", err)))
		return
	}
	fmt.Printf("rpc server listen at:%v with TLS\n", listen.Addr())

	// rpc server connection tcp listener
	err = server.Serve(listen)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("failed to serve: %v", err)))
	}
}
