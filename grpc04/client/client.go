package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc04/proto"
	"log"
	"strconv"
)

const Address string = ":8000"

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("grpc.Dial error %v", err)
	}

	defer conn.Close()

	grpcClient := pb.NewStreamClientClient(conn)

	stream, err := grpcClient.RouteList(context.Background())
	if err != nil {
		log.Fatalf("Call RouteList error %v", err)
	}

	for n := 0; n < 5; n++ {
		err = stream.Send(&pb.StreamRequest{
			StreamData: "stream client rpc " + strconv.Itoa(n),
		})

		if err != nil {
			log.Fatalf("strem.Send() error %v", err)
		}
	}

	req, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("stream.CloseAndRecv() error %v", err)
	}

	log.Println(req)
}
