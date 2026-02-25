package meta_ctx

import (
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
	"time"
)

func GetCtxByGinCtx(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	requestId := api.GenerateMsgIDFromContext(c)
	userId := api.GetUserIdFromContext(c)
	grpcCtx, cancel := GetCtxWithValueAndTimeout(timeout, userId, requestId)
	return grpcCtx, cancel
}

func GetCtxWithValue(userId int64, requestId string) context.Context {
	c := context.WithValue(context.Background(), consts.UserIdKey, userId)
	c = context.WithValue(c, consts.RequestIdKey, requestId)
	return c
}

func GetCtxWithValueAndTimeout(timeout time.Duration, userId int64, requestId string) (context.Context, context.CancelFunc) {
	c := GetCtxWithValue(userId, requestId)
	ctx, cancel := context.WithTimeout(c, timeout)
	return ctx, cancel
}

func NewGrpcMetaCtx(ctx context.Context) context.Context {
	md := metadata.MD{}
	value := ctx.Value(consts.RequestIdKey)
	if value != nil {
		if requestId, ok := value.(string); ok {
			SetRequestId(&md, requestId)
		}
	}
	value = ctx.Value(consts.UserIdKey)
	if value != nil {
		if userId, ok := value.(int64); ok {
			SetUserId(&md, userId)
		}
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func SetUserId(md *metadata.MD, userId int64) {
	md.Set(consts.UserIdKey, strconv.FormatInt(userId, 10))
}
func SetRequestId(md *metadata.MD, requestId string) {
	md.Set(consts.RequestIdKey, requestId)
}

func GetRequestIdFromCtx(ctx context.Context) string {
	if v := ctx.Value(strings.ToLower(consts.RequestIdKey)); v != nil {
		switch v.(type) {
		case []string:
			ids := v.([]string)
			if len(ids) > 0 {
				return ids[0]
			}
		case string:
			return v.(string)
		default:
			return ""
		}
	}
	return ""
}
func GetUserIdFromCtx(ctx context.Context) int64 {
	if v := ctx.Value(strings.ToLower(consts.UserIdKey)); v != nil {
		switch v.(type) {
		case []string:
			ids := v.([]string)
			if len(ids) > 0 {
				id, _ := strconv.ParseInt(ids[0], 10, 64)
				return id
			}
		case string:
			id, _ := strconv.ParseInt(v.(string), 10, 64)
			return id
		case int64:
			return v.(int64)
		default:
			return 0
		}
	}
	return 0
}
