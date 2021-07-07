package main

import (
	"context"
	"fmt"
	"grpcdemo-go/protobuf"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	address := "localhost:9000"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("connect fail %v\n", err)
	}
	defer conn.Close()

	c := protobuf.NewHelloServiceClient(conn)

	requestMessage := &protobuf.Massage{Request: "world "}
	responseMessage, err := c.Hello(context.Background(), requestMessage)
	if err != nil {
		fmt.Printf("hello err:%v\n ", err)
	}

	fmt.Printf("%v\n", responseMessage.String())

}
