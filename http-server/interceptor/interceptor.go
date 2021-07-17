package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// Interceptor 客户端拦截器
func Interceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, resp, cc, opts...)
	fmt.Printf("method=%s req=%v resp=%v duration=%s error=%v\n", method, req, resp, time.Since(start), err)
	return err
}
