package interceptor

import (
	"biz-auto-api/pkg/consts"
	grpcservice "biz-auto-api/pkg/grpc/meta_ctx"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/tools"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strings"
)

func UnaryValidator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := Validator(ctx, req); err != nil {

			return nil, err
		}
		return handler(ctx, req)
	}
}

func StreamValidator() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO 待完成关于流式接口的参数校验
		return handler(srv, stream)
	}
}

func Validator(ctx context.Context, req any) error {
	v := validator.New()
	err := v.Struct(req)
	if err == nil {
		return nil
	}
	msg := "参数不合法: "
	es := err.(validator.ValidationErrors)
	var errFields []string
	for _, e := range es {
		m := ""
		if len(e.Param()) > 0 && !strings.Contains(e.ActualTag(), e.Param()) {
			m = fmt.Sprintf("%v(%v=%v)", e.StructField(), e.ActualTag(), e.Param())
		} else {
			m = fmt.Sprintf("%v(%v)", e.StructField(), e.ActualTag())
		}
		if !tools.InSlice(m, errFields) {
			errFields = append(errFields, m)
		}
	}
	// 拼接不符合要求的字段
	msg += strings.Join(errFields, ",")
	err = errors.New(msg)
	requestId := grpcservice.GetRequestIdFromCtx(ctx)
	logger.GetLogger().WithField(strings.ToLower(consts.RequestIdKey), requestId).Errorf("%v", err)
	return err
}
