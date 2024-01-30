package userServer

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"sso_3.0/internal/services/user"
	api "sso_3.0/proto/gen"
)

type serverApi struct {
	userService *user.Service
	api.UnimplementedUserApiServer
	log *slog.Logger
}

func RegisterServer(grpcServer *grpc.Server, userService *user.Service, log *slog.Logger) {
	api.RegisterUserApiServer(grpcServer, &serverApi{userService: userService, log: log})
}

func (s *serverApi) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	//username := req.GetUsername()
	//user := s.userService.
	return &api.CreateUserResponse{}, nil
}
