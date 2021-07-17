package main

import (
	"context"
	"fmt"
	"grpcdemo-go/http-server/interceptor"
	"grpcdemo-go/protobuf"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type authCredential struct{}

func (*authCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": "123456", "secret_key": "secret_key"}, nil
}

func (*authCredential) RequireTransportSecurity() bool {
	return true
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endPoint := "127.0.0.1:9000"
	mux := runtime.NewServeMux()

	opts := make([]grpc.DialOption, 0)

	// 客户端请求拦截
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor.Interceptor))

	// TLS认证
	creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "www.zhouxuwen.com")
	if err != nil {
		fmt.Printf("failed to generate credentials err:%v\n", err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// 自定义认证方式
	opts = append(opts, grpc.WithPerRPCCredentials(new(authCredential)))

	// HTTP转grpc
	err = protobuf.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, endPoint, opts)
	if err != nil {
		fmt.Printf("Register handler err:%v\n", err)
	}

	fmt.Printf("HTTP Listen on 8080\n")
	http.ListenAndServe(":8080", mux)

}
