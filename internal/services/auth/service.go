package authService

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log/slog"
	userModel "sso_3.0/internal/domain/user"
	appErrors "sso_3.0/internal/errors"
	"sso_3.0/internal/pkg/bcrypt"
	"sso_3.0/internal/pkg/jwt"
	"sso_3.0/internal/storage/postgres"
	"sso_3.0/internal/storage/postgres/user"
	"strings"
)

type Service struct {
	log         *slog.Logger
	userStorage *user.Storage
}

func New(log *slog.Logger, storage *postgres.Storage) *Service {
	return &Service{log: log, userStorage: storage.UserStorage}
}

func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
	op := "service.auth.Register"
	log := s.log.With("op", op)

	hash, err := bcrypt.HashPassword(password)
	if err != nil {
		log.Error("Error on Hashing Password")
		return "", err
	}
	user, err := s.userStorage.Register(ctx, email, hash)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewToken(user)

	return token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	op := "service.auth.Login"
	log := s.log.With("op", op)

	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(appErrors.ErrUserNotExists, err) {
			return "", appErrors.ErrInvalidCredentials
		}
		log.Error("Error: ", err)
		return "", err
	}

	err = bcrypt.CheckPasswordHash(password, user.Hash)
	if err != nil {
		if errors.Is(appErrors.ErrPasswordIncorrect, err) {
			return "", appErrors.ErrInvalidCredentials
		}
		log.Error("Error: ", err)
		return "", err
	}

	token, err := jwt.NewToken(user)

	return token, nil
}

func (s *Service) ValidateAuth(ctx context.Context) (error, *userModel.Model) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return appErrors.NoTokenSent, nil
	}

	tokens := md.Get("authorization")

	if len(tokens) == 0 {
		return appErrors.NoTokenSent, nil
	}

	token := tokens[0]

	jwtToken := strings.TrimPrefix(token, "Bearer ")
	err, uid := s.ValidateToken(ctx, jwtToken)
	if err != nil || uid == "" {
		fmt.Println(err)
		return err, nil
	}

	user, err := s.userStorage.GetUserById(ctx, uid)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	return nil, user
}

func (s *Service) ValidateToken(ctx context.Context, token string) (error, string) {
	uid, err := jwt.CheckToken(token)
	if err != nil {
		return err, ""
	}

	if err != nil {
		return err, ""
	}

	return nil, uid
}

func (s *Service) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if !s.checkIfRoutePrivate(info.FullMethod) {
		return handler(ctx, req)
	}

	err, user := s.ValidateAuth(ctx)

	if err != nil {
		fmt.Printf("error: %e", err)
		if errors.Is(appErrors.NoTokenSent, err) {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, grpc.Errorf(codes.Unauthenticated, "Auth error")
	}

	ctx = context.WithValue(ctx, "uid", user.Id)
	ctx = context.WithValue(ctx, "email", user.Email)

	return handler(ctx, req)
}

func (s *Service) checkIfRoutePrivate(route string) bool {
	public := []string{
		"/api.AuthApi/Login",
		"/api.AuthApi/Register",
	}

	for _, item := range public {
		if item == route {
			return false
		}
	}

	return true
}

func (s *Service) GetUserFromCTX(ctx context.Context) *userModel.Model {
	email := ctx.Value("email").(string)
	uid := ctx.Value("uid").(string)

	return &userModel.Model{Id: uid, Email: email}

}
