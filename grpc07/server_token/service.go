package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "grpc_test/grpc01/proto"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	if err := CheckToken(ctx); err != nil {
		return nil, err
	}

	res := pb.SimpleResponse{Code: 200, Value: "ok " + req.Data}
	return &res, nil
}

func CheckToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "获取token失败")
	}

	var (
		appID     string
		appSecret string
	)

	if value, ok := md["app_id"]; ok {
		appID = value[0]
	}

	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appID != "grpc_token" || appSecret != "123456" {
		return status.Errorf(codes.Unauthenticated, "Token无效：app_id=%s，app_secret=%s", appID, appSecret)
	}

	return nil
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.listen error %v", err)
	}

	log.Println(Address + " net.listen...")

	// grpcServer := grpc.NewServer()
	// pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 启用拦截器 UnaryServerInterceptor 普通拦截器 StreamInterceptor流式拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = CheckToken(ctx)
		if err != nil {
			return
		}

		return handler(ctx, req)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("grpcServer.Serve() %v", err)
	}
}
