package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "grpc_test/grpc01/proto"
	"log"
	"net"
)

type SimpleService struct {
}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	tokenInfo := ctx.Value("token_info").(TokenInfo)
	res := pb.SimpleResponse{
		Code:  200,
		Value: "grpc auth token ok: id=" + tokenInfo.ID,
	}

	return &res, nil
}

type TokenInfo struct {
	ID string
}

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := ParseToken(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, " %v", err)
	}

	// 使用context.WithValue添加了值后，可以用Value(key)方法获取值
	newCtx := context.WithValue(ctx, "token_info", tokenInfo)

	return newCtx, nil
}

// ParseToken 解析token
func ParseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo

	sDec, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return tokenInfo, errors.New(fmt.Sprintf("解析Token: bearer %s 失败：%v", token, err))
	}

	if len(sDec) > 0 {
		tokenInfo.ID = string(sDec)
		return tokenInfo, nil
	}

	return tokenInfo, errors.New("Token无效：bearer " + token)
}

const (
	Address string = ":8000"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen error %v", err)
	}

	log.Println(Address + " net.Listen...")

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(AuthInterceptor),
		),
	))

	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("grpcServer.Serve() error %v", err)
	}
}
