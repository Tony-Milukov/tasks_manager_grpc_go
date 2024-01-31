package userServer

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	appErrors "sso_3.0/internal/errors"
	"sso_3.0/internal/services/user"
	api "sso_3.0/proto/gen"
)

type serverApi struct {
	userService *userService.Service
	api.UnimplementedUserApiServer
	log *slog.Logger
}

func RegisterServer(grpcServer *grpc.Server, userService *userService.Service, log *slog.Logger) {
	api.RegisterUserApiServer(grpcServer, &serverApi{userService: userService, log: log})
}

func (s *serverApi) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	email := req.GetEmail()
	pwd := req.GetPassword()
	token, err := s.userService.Register(ctx, email, pwd)

	if err != nil {
		if errors.Is(appErrors.ErrUserExists, err) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return &api.RegisterResponse{Token: token}, nil
}
func (s *serverApi) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	email := req.GetEmail()
	pwd := req.GetPassword()
	token, err := s.userService.Login(ctx, email, pwd)

	if err != nil {
		if errors.Is(appErrors.ErrInvalidCredentials, err) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return &api.LoginResponse{Token: token}, nil
}
