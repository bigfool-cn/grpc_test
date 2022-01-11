package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc_test/grpc01/proto"
	"log"
	"net"
)

type SimpleServer struct {
}

func (s *SimpleServer) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := &pb.SimpleResponse{Code: 200, Value: "ok " + req.Data}
	return res, nil
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)

	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	creds, err := credentials.NewServerTLSFromFile("../server.pem", "../server.key")

	if err != nil {
		log.Fatalf("Failed to generate credentials error %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterSimpleServer(grpcServer, &SimpleServer{})

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalf("grpcServer.Serve() error %v", err)
	}
}
