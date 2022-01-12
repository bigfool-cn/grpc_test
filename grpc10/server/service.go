package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	pb "grpc_test/grpc10/proto"
	"log"
	"net"
	"net/http"
	"strings"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.InnerMessage) (*pb.OuterMessage, error) {
	res := pb.OuterMessage{ImportantString: "ok", Inner: req}

	return &res, nil
}

func ProvideHTTP(endpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()

	// 新建gwmux，它是grpc-gateway的请求复用器。它将http请求与模式匹配，并调用相应的处理程序。
	gwmux := runtime.NewServeMux()
	// 将服务的http处理程序注册到gwmux。处理程序通过endpoint转发请求到grpc端点
	err := pb.RegisterSimpleHandlerFromEndpoint(ctx, gwmux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Register Endpoint err: %v", err)
	}

	// 新建mux，它是http的请求复用器
	mux := http.NewServeMux()
	// 注册gwmux
	mux.Handle("/", gwmux)
	log.Println(endpoint + " HTTP.Listening...")
	return &http.Server{
		Addr:    endpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
}

// grpcHandlerFunc 根据不同的请求重定向到指定的Handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen() error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_validator.UnaryServerInterceptor())))

	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 使用gateway把grpcServer转成httpServer
	httpServer := ProvideHTTP(Address, grpcServer)
	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatalf("httpServer.Serve() error %v", err)
	}
}
