package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc_test/grpc01/proto"
	"log"
)

const Address string = ":8000"

func main() {
	creds, err := credentials.NewClientTLSFromFile("../server.pem", "grpc")
	if err != nil {
		log.Fatalf("Failed create TLS credentials error %v", err)
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("grpc.Dial error %v", err)
	}

	defer conn.Close()

	grpcClient := pb.NewSimpleClient(conn)

	req := &pb.SimpleRequest{Data: "grpc tls"}
	res, err := grpcClient.Route(context.Background(), req)
	if err != nil {
		log.Fatalf("Call Route() error %v", err)
	}
	log.Println(res)
}
