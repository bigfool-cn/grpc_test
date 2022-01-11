package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
	pb "grpc_test/grpc01/proto"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}

	return &res, nil
}

func ZapInterceptor() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize za logger: %v", err)
	}
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}

func ZapInterceptorFile() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/debug.log",
		MaxSize:   1024,
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, zap.NewAtomicLevel())

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
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
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(grpc_zap.UnaryServerInterceptor(ZapInterceptorFile())),
		),
	)
	// 注册服务
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})
	// 调用grpc Serve()阻塞等待，直到进程退出或调用Stop()
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("gprc.Serve error %v", err)
	}
}
