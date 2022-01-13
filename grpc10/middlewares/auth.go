package middlewares

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TokenInfo struct {
	ID string
}

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, " %v", err)
	}

	// 使用context.WithValue添加了值后，可以用Value(key)方法获取值
	newCtx := context.WithValue(ctx, "token_info", tokenInfo)

	return newCtx, nil
}

// parseToken 解析token
func parseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo

	sDec, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return tokenInfo, errors.New(fmt.Sprintf("解析Token: bearer %s 失败：%v", token, err))
	}

	if string(sDec) == "grpc.auth.token" {
		tokenInfo.ID = string(sDec)
		return tokenInfo, nil
	}

	return tokenInfo, errors.New("Token无效：bearer " + token)
}
