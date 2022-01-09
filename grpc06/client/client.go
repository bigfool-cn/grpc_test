package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpc_test/grpc01/proto"
	"log"
	time2 "time"
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

	// 设置超时时间为3s
	clientDealTime := time2.Now().Add(3 * time2.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDealTime)
	defer cancel()

	// 调用服务方法
	res, err := grpcClient.Route(ctx, &req)
	if err != nil {
		//获取错误状态
		statu, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if statu.Code() == codes.DeadlineExceeded {
				log.Fatalln("Route timeout!")
			}
		}
		log.Fatalf("Call route error %v", err)
	}

	log.Println(res.Value)
}
