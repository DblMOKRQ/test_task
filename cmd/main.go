package main

import (
	"context"
	"fmt"

	"github.com/DblMOKRQ/test_task/internal/config"
	"github.com/DblMOKRQ/test_task/internal/repository"
	"github.com/DblMOKRQ/test_task/internal/service"
	"github.com/DblMOKRQ/test_task/internal/storage"
	rout "github.com/DblMOKRQ/test_task/internal/transport/rest"
	"github.com/DblMOKRQ/test_task/internal/transport/rest/handlers"
	"github.com/DblMOKRQ/test_task/pkg/logger"
	"go.uber.org/zap"
)

func main() {

	ctx := context.Background()

	log := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	storage, err := storage.New(ctx, cfg)
	if err != nil {
		log.Fatal("failed to create storage", zap.Error(err))
	}

	repo := repository.NewRepository(storage)

	service := service.NewService(repo, log, ctx)

	handler := handlers.NewHandlers(service)

	r := rout.NewRouter(handler)

	if err := r.Run(fmt.Sprintf("%s:%s", cfg.RestHost, cfg.RestPort)); err != nil {
		log.Fatal("failed to run server", zap.Error(err))
	}

}

// TODO:
// 2. cfg
