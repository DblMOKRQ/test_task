package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	systemLog "log"
	"net/http"

	"testtask/internal/config"
	postgres "testtask/internal/repository"
	"testtask/internal/service"
	"testtask/internal/transport/http/handler"
	"testtask/internal/transport/http/router"
	"testtask/pkg/logger"
)

func main() {

	ctx := context.Background()

	cfg := config.MustLoad()
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := log.Sync()
		if err != nil {
			systemLog.Printf("failed to sync logger: %v", err)
		}
	}()

	storeRepo, err := postgres.NewStore(ctx, cfg.UserRepo, cfg.PasswordRepo, cfg.HostRepo, cfg.PortRepo, cfg.DBName, cfg.SSLMode, log)
	if err != nil {
		log.Error("Failed to initialized to postgres", zap.Error(err))
		return
	}

	walletSrv := service.NewWalletService(storeRepo, log)

	handl := handler.NewHandler(*walletSrv)
	rout := router.NewRouter(handl, cfg.LogLevel, log)
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: rout.GetEngine(),
	}
	log.Info("Starting server", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Failed to listen and server", zap.Error(err))
		return
	}
}
