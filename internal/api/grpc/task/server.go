package taskServer

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	appErrors "sso_3.0/internal/errors"
	authService "sso_3.0/internal/services/auth"
	"sso_3.0/internal/services/tasks"
	api "sso_3.0/proto/gen"
)

type serverApi struct {
	authService *authService.Service
	taskService *tasks.Service
	log         *slog.Logger
	api.UnimplementedTaskApiServer
}

func RegisterServer(grpcServer *grpc.Server, authService *authService.Service, taskService *tasks.Service, log *slog.Logger) {
	api.RegisterTaskApiServer(grpcServer, &serverApi{authService: authService, taskService: taskService, log: log})
}

func (s *serverApi) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	title := req.GetTitle()
	description := req.GetDescription()
	statusId := req.GetStatusId()
	due := req.GetDue().AsTime()
	user := s.authService.GetUserFromCTX(ctx)

	task, err := s.taskService.CreateTask(ctx, title, description, user.Id, int(statusId), due)
	fmt.Println(task)
	if err != nil {
		if errors.Is(appErrors.ErrStatusUndefined, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	status := &api.Status{
		Title:       task.Status.Title,
		Description: task.Status.Description,
		Id:          int64(task.Status.Id),
	}

	return &api.CreateTaskResponse{
		Description: task.Description,
		Title:       task.Title,
		Id:          int64(task.Id),
		CreatorId:   task.CreatorId,
		Due:         timestamppb.New(task.Due),
		StatusId:    status.Id,
		Status:      status,
	}, err
}
func (s *serverApi) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	taskId := req.GetTaskId()

	err := s.taskService.DeleteTask(ctx, int(taskId))
	if err != nil {
		if errors.Is(appErrors.NothingToDelete, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	return &api.DeleteTaskResponse{Status: "Success"}, nil
}
func (s *serverApi) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTask not implemented")
}
func (s *serverApi) CreateStatus(ctx context.Context, req *api.CreateStatusRequest) (*api.CreateStatusResponse, error) {
	title := req.GetTitle()
	description := req.GetDescription()

	statusRes, err := s.taskService.CreateStatus(ctx, title, description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	return &api.CreateStatusResponse{Description: statusRes.Description, Title: statusRes.Title, Id: int64(statusRes.Id)}, err
}
func (s *serverApi) DeleteStatus(ctx context.Context, req *api.DeleteStatusRequest) (*api.DeleteStatusResponse, error) {
	statusId := req.GetStatusId()

	err := s.taskService.DeleteStatus(ctx, int(statusId))
	if err != nil {
		if errors.Is(appErrors.NothingToDelete, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	return &api.DeleteStatusResponse{Status: "Success"}, nil
}
func (s *serverApi) UpdateStatus(ctx context.Context, req *api.UpdateStatusRequest) (*api.UpdateStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
