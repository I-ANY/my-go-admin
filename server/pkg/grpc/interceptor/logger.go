package interceptor

import (
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/grpc/meta_ctx"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sirupsen/logrus"
	"strings"
)

func Logger(l *logrus.Logger) logging.Logger {
	return logging.LoggerFunc(func(c context.Context, lvl logging.Level, msg string, fields ...any) {

		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		for i.Next() {
			k, v := i.At()
			f[k] = v
		}
		requestId := meta_ctx.GetRequestIdFromCtx(c)
		f[strings.ToLower(consts.RequestIdKey)] = requestId
		l := l.WithFields(f)
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			l.Info(msg)
		}
	})
}
