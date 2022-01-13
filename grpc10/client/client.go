package main

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc"
	pb "grpc_test/grpc10/proto"
	"log"
)

const Address string = ":8000"
const headerAuthization string = "Authorization"

type Token struct {
	Value string
}

func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthization: t.Value}, nil
}

func (t *Token) RequireTransportSecurity() bool {
	return false
}

func main() {
	token := Token{
		Value: "bearer " + base64.StdEncoding.EncodeToString([]byte("grpc.auth.token")),
	}

	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&token))
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
