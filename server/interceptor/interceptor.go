package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, err := Auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

// ClientInterceptor 请求拦截器
func ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, resp, cc, opts...)
	fmt.Printf("method=%s req=%v resp=%v duration=%s error=%v\n", method, req, resp, time.Since(start), err)
	return err
}
