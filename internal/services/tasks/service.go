package tasks

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
	"sso_3.0/internal/domain/models"
	"sso_3.0/internal/domain/user"
	appErrors "sso_3.0/internal/errors"
	"sso_3.0/internal/storage/postgres"
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

	fmt.Printf("\n\nCreatorid: %s\n\n", creatorId)
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
