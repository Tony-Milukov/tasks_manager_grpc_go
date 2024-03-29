package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authServer "sso_3.0/internal/api/grpc/auth"
	taskServer "sso_3.0/internal/api/grpc/task"
	configParser "sso_3.0/internal/config"
	authService "sso_3.0/internal/services/auth"
	"sso_3.0/internal/services/tasks"
	"strconv"
)

type App struct {
	port       int
	grpcServer *grpc.Server
	log        *slog.Logger
}

func New(logger *slog.Logger, cfg *configParser.Config, authService *authService.Service, taskService *tasks.Service) (*App, error) {
	const op = "app.grpc.New"
	log := logger.With("op", op)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authService.AuthInterceptor))
	authServer.RegisterServer(grpcServer, authService, log)
	taskServer.RegisterServer(grpcServer, authService, taskService, log)

	log.Info("Servers successfully registered")
	port, err := strconv.Atoi(cfg.GrpcPort)

	if err != nil {
		return nil, err
	}
	return &App{grpcServer: grpcServer, port: int(port), log: log}, nil
}

// run it creates tcp listener and starts grpc server
func (s *App) run() error {
	op := "grpc.app.RUN"
	//setup getLogger
	log := s.log.With("op", op)

	//starting TCP listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		return err
	}

	//register grpc server with tcp listener
	err = s.grpcServer.Serve(l)
	if err != nil {
		return err
	}

	log.Info("Successfully Started GRPC api", "port", s.port)
	return nil
}

// MustRun Runs the application, if there is an errors it panics
func (app *App) MustRun() {
	if err := app.run(); err != nil {
		panic(err)
	}
}
