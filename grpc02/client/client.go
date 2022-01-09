package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_test/grpc01/proto"
	"log"
)

const Address string = ":8000"

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial error %v", err)
	}

	defer conn.Close()

	// 建立grpc连接
	grpcClient := pb.NewSimpleClient(conn)

	req := pb.SimpleRequest{
		Data: "grpc",
	}

	// 调用服务方法
	// 同时传入一个上下文context，方便控制grpc行为，比如超时/取消操作
	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call route error %v", err)
	}

	log.Println(res)
}
