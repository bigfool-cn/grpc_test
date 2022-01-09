package main

import (
	"google.golang.org/grpc"
	pb "grpc_test/grpc03/proto"
	"log"
	"net"
	"strconv"
)

type StreamSimpleService struct {
}

func (s *StreamSimpleService) ListValue(req *pb.SimpleRequest, svc pb.StreamSimple_ListValueServer) error {
	for n := 0; n < 5; n++ {
		err := svc.Send(&pb.SimpleResponse{
			StreamValue: req.Data + strconv.Itoa(n),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

const (
	Address string = ":8000"
	Netwotk string = "tcp"
)

func main() {
	listen, err := net.Listen(Netwotk, Address)

	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer()
	pb.RegisterStreamSimpleServer(grpcServer, &StreamSimpleService{})

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalf("grpcServer.Serve error %v", err)
	}
}
