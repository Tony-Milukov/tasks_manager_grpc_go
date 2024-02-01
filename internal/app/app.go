package app

import (
	"log/slog"
	"sso_3.0/internal/app/grpc"
	configParser "sso_3.0/internal/config"
	"sso_3.0/internal/services/auth"
	"sso_3.0/internal/services/tasks"
	"sso_3.0/internal/storage/postgres"
)

type App struct {
	GrpcServer *grpc.App
}

// New It creates new object of App
func New(cfg *configParser.Config, log *slog.Logger) (*App, error) {
	//create all storages
	storage, err := postgres.New(cfg, log)

	if err != nil {
		return nil, err
	}

	//crate services
	taskService := tasks.New(log, storage)
	authService := authService.New(log, storage)

	grpcServer, err := grpc.New(log, cfg, authService, taskService)

	if err != nil {
		return nil, err
	}

	return &App{
		grpcServer,
	}, nil
}
