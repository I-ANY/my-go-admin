package interceptor

import (
	"biz-auto-api/pkg/consts"
	grpcservice "biz-auto-api/pkg/grpc/meta_ctx"
	"biz-auto-api/pkg/logger"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func UnaryRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			err = errors.Errorf("%+v", r)
			requestId := grpcservice.GetRequestIdFromCtx(ctx)
			logger.GetLogger().WithField(consts.RequestIdKey, requestId).Errorf("%+v", err)
		}()
		return handler(ctx, req)
	}
}

func StreamRecover() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			err = errors.Errorf("%+v", r)
			requestId := grpcservice.GetRequestIdFromCtx(stream.Context())
			logger.GetLogger().WithField(consts.RequestIdKey, requestId).Errorf("%+v", err)
		}()
		return handler(srv, stream)
	}
}
