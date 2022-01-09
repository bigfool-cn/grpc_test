package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc03/proto"
	"io"
	"log"
)

const Address string = ":8000"

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("grpc.Dial error %v", err)
	}

	defer conn.Close()

	req := pb.SimpleRequest{
		Data: "stream grpc ",
	}

	grpClient := pb.NewStreamSimpleClient(conn)

	stream, err := grpClient.ListValue(context.Background(), &req)

	if err != nil {
		log.Fatalf("Call ListVale error %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ListValue get strem error %v", err)
		}

		log.Println(res.StreamValue)
	}
}
