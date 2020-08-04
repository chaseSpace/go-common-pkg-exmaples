package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
)

/*
在 gRPC 中，大类可分为两种 RPC 方法，与拦截器的对应关系是：

普通方法：一元拦截器（grpc.UnaryInterceptor）
流方法：流拦截器（grpc.StreamInterceptor）

下面介绍：grpc.UnaryInterceptor 签名如下
type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
参数说明：
ctx context.Context：请求上下文
req interface{}：RPC 方法的请求参数
info *UnaryServerInfo：RPC 方法的所有信息
handler UnaryHandler：RPC 方法本身
*/

/*
gRPC 本身居然只能设置一个拦截器, 我们不能直接把不同作用的逻辑写在一个拦截器里面
采用开源项目 go-grpc-middleware 实现分离
*/

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			// 当一个拦截器返回err时，下一个拦截器不会执行
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}
