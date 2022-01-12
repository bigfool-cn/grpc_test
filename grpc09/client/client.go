package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc10/proto"
	"log"
)

const Address string = ":8000"

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial() error %v", err)
	}

	defer conn.Close()

	grpcClient := pb.NewSimpleClient(conn)

	req := pb.InnerMessage{SomeFloat: 1, SomeInteger: 99}

	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("grpcClient.Route() error %v", err)
	}

	log.Println(res)
}
