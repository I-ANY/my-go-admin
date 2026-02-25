package clients

import (
	"biz-auto-api/pkg/grpc/meta_ctx"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func WithUnaryMetadata() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcCtx := meta_ctx.NewGrpcMetaCtx(ctx)
		return invoker(grpcCtx, method, req, reply, cc, opts...)
	}
}

func WithStreamMetadata() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		grpcCtx := meta_ctx.NewGrpcMetaCtx(ctx)
		return streamer(grpcCtx, desc, cc, method, opts...)
	}
}

func GetClientOptions(t time.Duration) []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(t)),
		grpc.WithChainUnaryInterceptor(WithUnaryMetadata()),
		grpc.WithStreamInterceptor(WithStreamMetadata()),
		// 添加更多配置
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024), // 50MB
			grpc.MaxCallSendMsgSize(50*1024*1024), // 50MB
		),
	}
	return opts
}
