package logger

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func LoggerMiddleware() {

}

// func LoggerInterceptor(ctx context.Context,
// 	req interface{},
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (interface{}, error) {
// 	logger := NewLogger()

// 	logger.Info("request", zap.Any("Method", info.FullMethod),
// 		zap.Any("request", req),
// 	)

// 	return handler(ctx, req)
// }
