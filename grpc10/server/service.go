package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
	"grpc_test/grpc10/middlewares"
	pb "grpc_test/grpc10/proto"
	"grpc_test/grpc10/server/gateway"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.InnerMessage) (*pb.OuterMessage, error) {
	res := pb.OuterMessage{ImportantString: "ok", Inner: req}

	return &res, nil
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen() error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				grpc_auth.UnaryServerInterceptor(middlewares.AuthInterceptor),
			),
		),
	)

	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 使用gateway把grpcServer转成httpServer
	httpServer := gateway.ProvideHTTP(Address, grpcServer)
	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatalf("httpServer.Serve() error %v", err)
	}
}
