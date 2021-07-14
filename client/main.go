package main

import (
	"context"
	"fmt"
	"grpcdemo-go/client/interceptor"
	"grpcdemo-go/protobuf"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type authCredential struct{}

func (*authCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": "123456", "secret_key": "key"}, nil
}

func (*authCredential) RequireTransportSecurity() bool {
	return true
}

func main() {
	opts := make([]grpc.DialOption, 0)

	// TLS认证
	creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "www.zhouxuwen.com")
	if err != nil {
		fmt.Printf("failed to generate credentials err:%v\n", err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// 自定义认证方式
	// opts = append(opts, grpc.WithPerRPCCredentials(new(authCredential)))

	// 客户端请求拦截
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor.Interceptor))

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:9000", opts...)
	if err != nil {
		fmt.Printf("connect fail %v\n", err)
	}
	defer conn.Close()

	c := protobuf.NewHelloServiceClient(conn)

	// 认证参数
	md := metadata.Pairs("app_id", "123456", "secret_key", "secret_key")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	requestMessage := &protobuf.Massage{Request: "world"}
	responseMessage, err := c.Hello(ctx, requestMessage)
	if err != nil {
		fmt.Printf("hello err:%v\n ", err)
		return
	}

	fmt.Printf("%v\n", responseMessage.String())

}
