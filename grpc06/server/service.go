package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpc_test/grpc01/proto"
	"log"
	"net"
	"runtime"
	time2 "time"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	data := make(chan *pb.SimpleResponse, 1)
	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
}

func handle(ctx context.Context, req *pb.SimpleRequest, data chan<- *pb.SimpleResponse) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		runtime.Goexit()
	case <-time2.After(4 * time2.Second):
		res := pb.SimpleResponse{
			Code:  200,
			Value: "hello " + req.Data,
		}

		data <- &res
	}
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	// 创建grpc服务器实例
	grpcServer := grpc.NewServer()
	// 注册服务
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})
	// 调用grpc Serve()阻塞等待，直到进程退出或调用Stop()
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("gprc.Serve error %v", err)
	}
}
