package main

import (
	"errors"
	"fmt"
	"grpcdemo-go/protobuf"
	"grpcdemo-go/server/service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("failed to listen: %v", err)))
		return
	}
	fmt.Printf("rpc server listen at:%v\n", listen.Addr())

	server := grpc.NewServer()

	protobuf.RegisterHelloServiceServer(server, &service.HelloService)
	err = server.Serve(listen)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("failed to serve: %v", err)))
	}
}
