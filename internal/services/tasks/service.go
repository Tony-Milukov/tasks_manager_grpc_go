package tasks

import (
	"context"
	"log/slog"
	"sso_3.0/internal/domain/models"
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
func (s *Service) DeleteTask(ctx context.Context, id int) error {
	err := s.storage.TaskStorage.DeleteTask(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
func (s *Service) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*models.Task, error) {
	return nil, nil
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

func (s *Service) UpdateStatus(ctx context.Context, req *api.UpdateStatusRequest) (*models.Status, error) {
	return nil, nil
}
