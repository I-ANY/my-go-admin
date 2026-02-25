package clients

import (
	"biz-auto-api/pkg/grpc/pb/common"
	"biz-auto-api/pkg/grpc/pb/vpn"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

var (
	_vpnClient vpn.VpnServiceClient
)

func SetupVpnGrpcClient(addr string, timeout time.Duration) error {
	conn, err := grpc.NewClient(addr, GetClientOptions(timeout)...)
	if err != nil {
		return errors.WithStack(err)
	}

	// 等待连接就绪
	c := vpn.NewVpnServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = c.Ping(ctx, &common.EmptyReq{})
	if err != nil {
		return errors.WithStack(err)
	}
	_vpnClient = c
	return nil
}
func GetVpnGrpcClient() vpn.VpnServiceClient {
	return _vpnClient
}
