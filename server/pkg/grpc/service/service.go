package service

import (
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/grpc/meta_ctx"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type Service struct {
	db     *gorm.DB
	logger *logrus.Logger
	redis  *redis.Client
	ck     *clickhouse.CK
}
type Option func(*Service)

func (s *Service) WithOption(options ...Option) {
	for _, o := range options {
		o(s)
	}
}

func WithDB(db *gorm.DB) Option {
	return func(s *Service) {
		s.db = db
	}
}
func WithLogger(logger *logrus.Logger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithRedis(redis *redis.Client) Option {
	return func(s *Service) {
		s.redis = redis
	}
}

func WithClickhouse(ck *clickhouse.CK) Option {
	return func(s *Service) {
		s.ck = ck
	}
}

func (s *Service) GetUserId(ctx context.Context) int64 {
	return meta_ctx.GetUserIdFromCtx(ctx)
}

func (s *Service) GetRequestId(ctx context.Context) string {
	return meta_ctx.GetRequestIdFromCtx(ctx)
}

func (s *Service) wrapperLogger(ctx context.Context, log *logrus.Logger) *logrus.Entry {
	requestId := s.GetRequestId(ctx)
	return log.WithField(strings.ToLower(consts.RequestIdKey), requestId)
}

func (s *Service) wrapperDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	requestId := s.GetRequestId(ctx)
	c := context.WithValue(ctx, consts.RequestIdKey, requestId)
	return db.WithContext(c)
}

func (s *Service) GetDB(ctx context.Context) *gorm.DB {
	return s.wrapperDB(ctx, s.db)
}

func (s *Service) GetLogger(ctx context.Context) *logrus.Entry {
	return s.wrapperLogger(ctx, s.logger)
}

func (s *Service) GetRedis() *redis.Client {
	return s.redis
}
func (s *Service) GetClickhouse(ctx context.Context) *clickhouse.CK {
	requestId := s.GetRequestId(ctx)
	c := context.WithValue(ctx, consts.RequestIdKey, requestId)
	return s.ck.WithContext(c)
}
