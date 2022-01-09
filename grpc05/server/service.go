package main

import (
	"google.golang.org/grpc"
	pb "grpc_test/grpc05/proto"
	"io"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) RouteStream(srv pb.Stream_RouteStreamServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Println(req.Question)

		err = srv.Send(&pb.StreamResponse{
			Answer: "answer: " + req.Question,
		})

		if err != nil {
			return err
		}
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
		log.Fatalf("net.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer()

	pb.RegisterStreamServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("grpcServer.Serve() error %v", err)
	}
}
