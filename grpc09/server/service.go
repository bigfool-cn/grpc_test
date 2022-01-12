package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
	pb "grpc_test/grpc09/proto"
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_validator.UnaryServerInterceptor())))

	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalf("grpcServer.Serve() error %v", err)
	}
}
