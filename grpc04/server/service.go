package main

import (
	"google.golang.org/grpc"
	pb "grpc_test/grpc04/proto"
	"io"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	for {
		res, err := srv.Recv()

		if err == io.EOF {
			return srv.SendAndClose(&pb.SimpleReponse{Value: "ok"})
		}

		if err != nil {
			return err
		}

		log.Println(res.StreamData)
	}
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("new.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer()

	pb.RegisterStreamClientServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listen)

	if err == nil {
		log.Fatalf("grpc.Serve error %v", err)
	}
}
