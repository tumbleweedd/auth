package interceptors

import (
	"context"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

type LoggerInterceptor struct {
	logger logger.Logger
}

func NewLoggerInterceptor(logger logger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		logger: logger,
	}
}

func (l *LoggerInterceptor) UnaryLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		ctx = context.WithValue(ctx, "logger", l.logger)

		res, err := handler(ctx, req)

		duration := time.Since(startTime)
		statusCode := status.Code(err)

		l.logger.Info(
			"method", info.FullMethod,
			"duration", duration.String(),
			"status_code", statusCode.String(),
		)

		if err != nil {
			l.logger.Error(
				"method", info.FullMethod,
				"error", err.Error(),
			)
		}

		return res, err
	}
}
