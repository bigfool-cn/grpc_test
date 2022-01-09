package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc05/proto"
	"io"
	"log"
	"strconv"
)

const Address string = ":8000"

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial() error %v", err)
	}

	defer conn.Close()

	grpcClient := pb.NewStreamClient(conn)

	stream, err := grpcClient.RouteStream(context.Background())

	if err != nil {
		log.Fatalf("grpcClient.RouteStream() error %v", err)
	}

	for n := 0; n < 5; n++ {
		err = stream.Send(&pb.StreamResquest{
			Question: "question " + strconv.Itoa(n),
		})

		if err != nil {
			log.Fatalf("stream.Send() error %v", err)
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream.Send() error %v", err)
		}

		log.Println(req)
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("stream.CloseSend() error %v", err)
	}
}
