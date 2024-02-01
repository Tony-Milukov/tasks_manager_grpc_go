package task

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sso_3.0/internal/domain/models"
	appErrors "sso_3.0/internal/errors"
	api "sso_3.0/proto/gen"
	"time"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(db *sql.DB, log *slog.Logger) *Storage {
	return &Storage{db: db, log: log}
}

func (s *Storage) CreateTask(ctx context.Context, title, description, creatorId string, statusId int, due time.Time) (*models.Task, error) {
	var id int
	fmt.Println(due)
	status, err := s.GetStatusById(ctx, statusId)

	if err != nil {
		return nil, appErrors.ErrStatusUndefined
	}

	err = s.db.QueryRowContext(ctx, "INSERT INTO tasks (title, description, statusid, creatorId, due) VALUES ($1,$2,$3,$4, $5) RETURNING id", title, description, statusId, creatorId, due).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &models.Task{Title: title, CreatorId: creatorId, StatusId: statusId, Status: status, Id: id, Due: due, Description: description}, nil
}
func (s *Storage) DeleteTask(ctx context.Context, id int) error {
	execContext, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return appErrors.NothingToDelete
	}

	return nil
}
func (s *Storage) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*models.Task, error) {
	return nil, nil
}
func (s *Storage) CreateStatus(ctx context.Context, title, description string) (*models.Status, error) {
	var id int

	err := s.db.QueryRowContext(ctx, "INSERT INTO statuses (title, description) VALUES ($1,$2) RETURNING id", title, description).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &models.Status{
		Id:          id,
		Description: description,
		Title:       title,
	}, nil
}
func (s *Storage) DeleteStatus(ctx context.Context, id int) error {

	execContext, err := s.db.ExecContext(ctx, "DELETE FROM statuses WHERE id = $1", id)
	if err != nil {
		return err
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return appErrors.NothingToDelete
	}

	return nil
}

func (s *Storage) GetStatusById(ctx context.Context, id int) (*models.Status, error) {
	var title, description string
	err := s.db.QueryRowContext(ctx, "SELECT title, description FROM statuses WHERE id = $1", id).Scan(&title, &description)
	if err != nil {
		return nil, err
	}

	return &models.Status{Id: id, Title: title, Description: description}, nil
}

func (s *Storage) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	var title, description, creatorId string
	var statusId int
	err := s.db.QueryRowContext(ctx, "SELECT title, description, creatorId, statusId FROM tasks WHERE id = $1", id).Scan(&title, &description, &creatorId, &statusId)
	if err != nil {
		return nil, err
	}

	status, err := s.GetStatusById(ctx, statusId)
	if err != nil {
		return nil, err
	}

	return &models.Task{Title: title, CreatorId: creatorId, StatusId: statusId, Status: status}, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, req *api.UpdateStatusRequest) (*models.Status, error) {
	return nil, nil
}
