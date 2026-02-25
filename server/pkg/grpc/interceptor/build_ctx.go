package interceptor

import (
	"biz-auto-api/pkg/consts"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func BuildCtx(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md) == 0 {
		return ctx
	}
	var withValueCtx = ctx
	for key, value := range md {
		withValueCtx = context.WithValue(withValueCtx, strings.ToLower(key), value)
	}
	// 判断是是否有requestId，没有就生成
	value := withValueCtx.Value(strings.ToLower(consts.RequestIdKey))
	if value != nil {
		if requestIds, ok := value.([]string); !ok || len(requestIds) == 0 || requestIds[0] == "" {
			withValueCtx = context.WithValue(withValueCtx, strings.ToLower(consts.RequestIdKey), []string{uuid.New().String()})
		}
	} else {
		withValueCtx = context.WithValue(withValueCtx, strings.ToLower(consts.RequestIdKey), []string{uuid.New().String()})
	}
	return withValueCtx
}

func UnaryBuildMetadataToCtx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		withValueCtx := BuildCtx(ctx)
		return handler(withValueCtx, req)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	newCtx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.newCtx
}
func StreamBuildMetadataToCtx() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		withValueCtx := BuildCtx(stream.Context())
		wrapped := &wrappedStream{
			ServerStream: stream,
			newCtx:       withValueCtx,
		}
		return handler(srv, wrapped)
	}
}
