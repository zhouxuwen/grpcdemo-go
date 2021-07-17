package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"grpcdemo-go/protobuf"
	"grpcdemo-go/server/interceptor"
	"grpcdemo-go/server/service"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func main() {
	endPoint := "localhost:9000"
	listen, err := net.Listen("tcp", endPoint)
	if err != nil {
		fmt.Printf("tcp listen err:%v\n", err)
		return
	}

	// rpc server
	opts := make([]grpc.ServerOption, 0)
	creds, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
	if err != nil {
		fmt.Printf("failed to generate credentials err:%v\n", err)
		return
	}
	// TLS认证
	opts = append(opts, grpc.Creds(creds))
	// 自定义认证
	opts = append(opts, grpc.UnaryInterceptor(interceptor.Interceptor))
	server := grpc.NewServer(opts...)
	protobuf.RegisterHelloServiceServer(server, &service.HelloService)

	// http gateway server
	gwOpts := make([]grpc.DialOption, 0)
	gwCreds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "www.zhouxuwen.com")
	if err != nil {
		fmt.Printf("failed to generate credentials err:%v\n", err)
		return
	}
	// tls认证
	gwOpts = append(gwOpts, grpc.WithTransportCredentials(gwCreds))
	// 自定义认证方式
	gwOpts = append(gwOpts, grpc.WithPerRPCCredentials(&interceptor.AuthCredential))
	// 请求拦截
	gwOpts = append(gwOpts, grpc.WithUnaryInterceptor(interceptor.ClientInterceptor))

	gwmux := runtime.NewServeMux()
	err = protobuf.RegisterHelloServiceHandlerFromEndpoint(context.Background(), gwmux, endPoint, gwOpts)
	if err != nil {
		fmt.Printf("register http gateway server fail. err: %v\n", err)
		return
	}

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	tlsConfig, err := getTLSConfig()
	if err != nil {
		fmt.Printf("get tls config err :%v\n", err)
		return
	}
	srv := &http.Server{
		Addr:      endPoint,
		Handler:   grpcHandlerFunc(server, mux),
		TLSConfig: tlsConfig,
	}

	fmt.Printf("gRPC and http listen at :%v\n", endPoint)
	err = srv.Serve(tls.NewListener(listen, tlsConfig))
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("failed to serve: %v", err)))
	}
}

func getTLSConfig() (*tls.Config, error) {
	cert, _ := ioutil.ReadFile("./keys/server.pem")
	key, _ := ioutil.ReadFile("./keys/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		fmt.Printf("TLS KeyPair err: %v\n", err)
		return nil, err
	}
	demoKeyPair = &pair
	return &tls.Config{
		Certificates: []tls.Certificate{*demoKeyPair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
	}, nil
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
