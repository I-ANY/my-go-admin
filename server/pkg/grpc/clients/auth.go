package clients

import (
	"biz-auto-api/pkg/grpc/pb/auth"
	"biz-auto-api/pkg/grpc/pb/common"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

var (
	_authClient auth.AuthServiceClient
)

func SetupAuthGrpcClient(addr string, timeout time.Duration) error {
	conn, err := grpc.NewClient(addr, GetClientOptions(timeout)...)
	if err != nil {
		return errors.WithStack(err)
	}
	// 等待连接就绪
	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	_, err = c.Ping(ctx, &common.EmptyReq{})
	if err != nil {
		return errors.WithStack(err)
	}
	_authClient = c
	return nil
}
func GetAuthGrpcClient() auth.AuthServiceClient {
	return _authClient
}
