package tasks

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
	"sso_3.0/internal/domain/models"
	"sso_3.0/internal/domain/user"
	appErrors "sso_3.0/internal/errors"
	"sso_3.0/internal/storage/postgres"
	api "sso_3.0/proto/gen"
	"time"
)

type Service struct {
	log     *slog.Logger
	storage *postgres.Storage
}

func New(log *slog.Logger, storage *postgres.Storage) *Service {
	return &Service{log: log, storage: storage}
}

func (s *Service) CreateTask(ctx context.Context, title, description, creatorId string, statusId int, due time.Time) (*models.Task, error) {

	task, err := s.storage.TaskStorage.CreateTask(ctx, title, description, creatorId, statusId, due)

	if err != nil {
		return nil, err
	}

	return task, nil
}
func (s *Service) DeleteTask(ctx context.Context, id int, currentUser *user.Model) error {
	task, err := s.storage.TaskStorage.GetTaskById(ctx, id)

	if err != nil {
		fmt.Println(err)
		if errors.Is(appErrors.ErrTaskNotExists, err) {
			return appErrors.NothingToDelete
		}
		return err
	}

	if task.CreatorId != currentUser.Id {
		return appErrors.ErrNoPermission
	}

	err = s.storage.TaskStorage.DeleteTask(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
func (s *Service) UpdateTask(ctx context.Context, title, description string, due time.Time, statusId, id int, completed *wrapperspb.BoolValue, user *user.Model) (*models.Task, error) {
	var status *models.Status = nil
	var err error

	// check if task exists
	task, err := s.GetTaskById(ctx, id)

	if err != nil {
		return nil, err
	}

	//check if Creator is owner of the task
	if user.Id != task.CreatorId {
		return nil, appErrors.ErrNoPermission
	}

	// check if status is to update
	if statusId != 0 {
		status, err = s.storage.TaskStorage.GetStatusById(ctx, statusId)
		if err != nil {
			return nil, err
		}
	}

	// update task
	task, err = s.storage.TaskStorage.UpdateTask(ctx, title, description, due, status, completed, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}
func (s *Service) CreateStatus(ctx context.Context, title, description string) (*models.Status, error) {
	status, err := s.storage.TaskStorage.CreateStatus(ctx, title, description)

	if err != nil {
		return nil, err
	}

	return status, nil
}

func (s *Service) DeleteStatus(ctx context.Context, id int) error {
	err := s.storage.TaskStorage.DeleteStatus(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, title, description string, statusId int) (*models.Status, error) {
	_, err := s.storage.TaskStorage.GetStatusById(ctx, statusId)

	if err != nil {
		return nil, err
	}

	status, err := s.storage.TaskStorage.UpdateStatus(ctx, title, description, statusId)
	if err != nil {
		return nil, err
	}
	return status, nil
}
func (s *Service) GetTaskById(ctx context.Context, taskId int) (*models.Task, error) {
	task, err := s.storage.TaskStorage.GetTaskById(ctx, taskId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetCreatedTasksByUID(ctx context.Context, userId string) ([]*models.Task, error) {
	tasks, err := s.storage.TaskStorage.GetCreatedTasksByUID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) GetProtoAssignees(assignees []*models.Assignee) []*api.TaskAssignee {
	var value []*api.TaskAssignee
	if assignees != nil {
		for _, assignee := range assignees {
			user := &api.User{
				Id:    assignee.User.Id,
				Email: assignee.User.Email,
			}
			value = append(value, &api.TaskAssignee{
				TaskId: int64(assignee.TaskId),
				Role:   assignee.Role,
				Id:     int64(assignee.Id),
				User:   user,
			})
		}
	}

	return value
}

func (s *Service) GetProtoStatus(status *models.Status) *api.Status {
	if status != nil {
		return &api.Status{
			Title:       status.Title,
			Description: status.Description,
			Id:          int64(status.Id),
		}
	}
	return nil
}
func (s *Service) GetProtoTask(task *models.Task) *api.Task {
	var status *api.Status
	var completed bool
	var assignees []*api.TaskAssignee

	if task == nil {
		return nil
	}

	if task.Status != nil {
		status = s.GetProtoStatus(task.Status)
	}

	if task.Completed != nil {
		completed = task.Completed.Value
	}

	if task.Assignees != nil {
		assignees = s.GetProtoAssignees(task.Assignees)
	}
	return &api.Task{
		Id:          int64(task.Id),
		Title:       task.Title,
		Description: task.Description,
		Due:         timestamppb.New(task.Due),
		Status:      status,
		CreatorId:   task.CreatorId,
		Completed:   completed,
		Assignees:   assignees,
	}
}

func (s *Service) GetProtoTasks(tasks []*models.Task) []*api.Task {
	var value []*api.Task

	for _, task := range tasks {
		protoTask := s.GetProtoTask(task)
		value = append(value, protoTask)
	}

	return value
}

func (s *Service) AssignTask(ctx context.Context, userId, role string, taskId int) (*models.Task, error) {
	fmt.Printf("TASK: %d  USER: %s ROLE %s", taskId, userId, role)
	task, err := s.storage.TaskStorage.AssignTask(ctx, userId, role, taskId)

	if err != nil {
		return nil, err
	}

	return task, nil
}
