package taskServer

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if err != nil {
		if errors.Is(appErrors.ErrStatusUndefined, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	taskProto := s.taskService.GetProtoTask(task)

	return &api.CreateTaskResponse{
		Id:          taskProto.Id,
		Title:       taskProto.Title,
		Description: taskProto.Description,
		CreatorId:   taskProto.CreatorId,
		Due:         taskProto.Due,
		Completed:   taskProto.Completed,
		Status:      taskProto.Status,
		Assignees:   taskProto.Assignees,
	}, nil
}
func (s *serverApi) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	taskId := req.GetTaskId()
	currentUser := s.authService.GetUserFromCTX(ctx)
	err := s.taskService.DeleteTask(ctx, int(taskId), currentUser)

	if err != nil {
		if errors.Is(appErrors.NothingToDelete, err) || errors.Is(appErrors.ErrNoPermission, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	return &api.DeleteTaskResponse{Status: "Success"}, nil
}
func (s *serverApi) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {
	title := req.GetTitle()
	description := req.GetDescription()
	statusId := req.GetStatusId()
	due := req.GetDue()
	completed := req.GetCompleted()
	id := req.GetTaskId()

	//get user from ctx -> from JWT
	user := s.authService.GetUserFromCTX(ctx)

	//update task
	task, err := s.taskService.UpdateTask(ctx, title, description, due.AsTime(), int(statusId), int(id), completed, user)

	// handle errors
	if err != nil {
		fmt.Println(err)
		if errors.Is(appErrors.ErrStatusUndefined, err) || errors.Is(appErrors.ErrTaskNotExists, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(appErrors.ErrNoPermission, err) {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	taskProto := s.taskService.GetProtoTask(task)
	return &api.UpdateTaskResponse{
		Id:          taskProto.Id,
		Title:       taskProto.Title,
		Description: taskProto.Description,
		CreatorId:   taskProto.CreatorId,
		Due:         taskProto.Due,
		Completed:   taskProto.Completed,
		Status:      taskProto.Status,
	}, nil
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
		fmt.Println(err)
		if errors.Is(appErrors.NothingToDelete, err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, appErrors.Internal.Error())
	}

	return &api.DeleteStatusResponse{Status: "Success"}, nil
}

func (s *serverApi) UpdateStatus(ctx context.Context, req *api.UpdateStatusRequest) (*api.UpdateStatusResponse, error) {
	statusId := req.GetStatusId()
	description := req.GetDescription()
	title := req.GetTitle()

	if title == "" && description == "" {
		return nil, status.Errorf(codes.InvalidArgument, appErrors.NoArguments.Error())
	}

	statusRes, err := s.taskService.UpdateStatus(ctx, title, description, int(statusId))

	if err != nil {
		if errors.Is(appErrors.ErrStatusUndefined, err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, appErrors.Internal.Error())
	}

	return &api.UpdateStatusResponse{
		Id:          int64(statusRes.Id),
		Title:       statusRes.Title,
		Description: statusRes.Description,
	}, nil

}

func (s *serverApi) GetCreatedTasksByUID(ctx context.Context, req *api.GetCreatedTasksByUIDRequest) (*api.GetCreatedTasksByUIDResponse, error) {
	userId := req.GetUserId()
	var tasks []*api.Task
	if userId == "" {
		user := s.authService.GetUserFromCTX(ctx)
		userId = user.Id
	}

	tasksRes, err := s.taskService.GetCreatedTasksByUID(ctx, userId)

	if err != nil {
		fmt.Println(err)
		if errors.Is(appErrors.ErrStatusUndefined, err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, appErrors.Internal.Error())
	}

	tasks = s.taskService.GetProtoTasks(tasksRes)

	return &api.GetCreatedTasksByUIDResponse{
		Tasks: tasks,
	}, nil
}

func (s *serverApi) AssignTask(ctx context.Context, req *api.AssignTaskRequest) (*api.AssignTaskResponse, error) {
	description := req.GetDescription()
	userId := req.GetUserId()
	taskId := req.GetTaskId()

	user := s.authService.GetUserFromCTX(ctx)

	task, err := s.taskService.GetTaskById(ctx, int(taskId))

	if err != nil {
		if errors.Is(appErrors.ErrTaskNotExists, err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, appErrors.Internal.Error())
	}
	if task.CreatorId != user.Id {
		return nil, status.Errorf(codes.PermissionDenied, appErrors.ErrNoPermission.Error())
	}

	if userId == "" {
		user := s.authService.GetUserFromCTX(ctx)
		userId = user.Id
	}

	task, err = s.taskService.AssignTask(ctx, userId, description, int(taskId))

	if err != nil {
		if errors.Is(appErrors.ErrStatusUndefined, err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(appErrors.TaskAlreadyAssigned, err) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, appErrors.Internal.Error())
	}

	protoTask := s.taskService.GetProtoTask(task)
	fmt.Println(protoTask)
	return &api.AssignTaskResponse{
		Task: protoTask,
	}, nil
}
