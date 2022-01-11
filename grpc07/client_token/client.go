package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc01/proto"
	"log"
)

type Token struct {
	AppID     string
	AppSecret string
}

func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppID, "app_secret": t.AppSecret}, nil
}

func (t *Token) RequireTransportSecurity() bool {
	return false
}

const Address string = ":8000"

func main() {
	token := Token{
		AppID:     "grpc_token",
		AppSecret: "123456",
	}
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		log.Fatalf("grpc.Dial() error %v", err)
	}

	defer conn.Close()

	grpcClient := pb.NewSimpleClient(conn)

	req := pb.SimpleRequest{Data: "grpc token"}

	res, err := grpcClient.Route(context.Background(), &req)

	if err != nil {
		log.Fatalf("Call grpcClient.Route() error %v", err)
	}

	log.Println(res)
}
