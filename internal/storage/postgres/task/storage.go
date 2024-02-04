package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
	"sso_3.0/internal/domain/models"
	appErrors "sso_3.0/internal/errors"
	"strings"
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
	status, err := s.GetStatusById(ctx, statusId)
	fmt.Printf("\n\nStatus START: %s\n\n", creatorId)
	fmt.Println(status)
	fmt.Printf("\n\nStatus END: %s\n\n", creatorId)

	if err != nil {
		return nil, appErrors.ErrStatusUndefined
	}

	err = s.db.QueryRowContext(ctx, "INSERT INTO tasks (title, description, statusid, creatorId, due) VALUES ($1,$2,$3,$4, $5) RETURNING id", title, description, statusId, creatorId, due).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &models.Task{Title: title, CreatorId: creatorId, Status: status, Id: id, Due: due, Description: description}, nil
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
func (s *Storage) UpdateTask(ctx context.Context, title, description string, due time.Time, status *models.Status, completed *wrapperspb.BoolValue, id int) (*models.Task, error) {
	var fields []string
	var values []interface{}
	var creatorId string
	var statusId sql.NullInt64
	var completedNull sql.NullBool
	key := 2
	values = append(values, id)

	//generate query

	if title != "" {
		fields = append(fields, fmt.Sprintf("title = $%d", key))
		values = append(values, title)
		key++
	}

	if description != "" {
		fields = append(fields, fmt.Sprintf("description = $%d", key))
		values = append(values, description)
		key++
	}

	if !due.IsZero() {
		fields = append(fields, fmt.Sprintf("due = $%d", key))
		values = append(values, due)
		key++
	}

	if completed != nil {
		fields = append(fields, fmt.Sprintf("completed = $%d", key))
		values = append(values, completed.Value)
		key++
		completedNull = sql.NullBool{Bool: completed.Value, Valid: true}
	}

	if status != nil {
		fields = append(fields, fmt.Sprintf("statusId = $%d", key))
		values = append(values, status.Id)
		key++
	}

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $1 RETURNING title, description, creatorId, statusId, due, completed", strings.Join(fields, ", "))

	//execute the update and get new values
	err := s.db.QueryRowContext(ctx, query, values...).Scan(&title, &description, &creatorId, &statusId, &due, &completedNull)

	if err != nil {
		// if nothing was updated
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrTaskNotExists
		}
		return nil, err
	}

	var statusResult *models.Status

	if statusId.Valid {
		statusResult, err = s.GetStatusById(ctx, int(statusId.Int64))
		if err != nil {
			return nil, err
		}
	}

	// get completed
	if completedNull.Valid {
		completed = wrapperspb.Bool(completedNull.Bool)
	}

	task := &models.Task{
		Id:          id,
		Title:       title,
		CreatorId:   creatorId,
		Completed:   completed,
		Description: description,
		Due:         due,
		Status:      statusResult,
	}

	return task, nil

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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE tasks SET statusId=null WHERE statusId=$1 ", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	execContext, err := tx.ExecContext(ctx, "DELETE FROM statuses WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		tx.Rollback()

		return err
	}

	if affected == 0 {
		tx.Rollback()
		return appErrors.NothingToDelete
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetStatusById(ctx context.Context, id int) (*models.Status, error) {
	var title, description string
	err := s.db.QueryRowContext(ctx, "SELECT title, description FROM statuses WHERE id = $1", id).Scan(&title, &description)
	if err != nil {
		fmt.Println(err)
		if errors.Is(sql.ErrNoRows, err) {
			return nil, appErrors.ErrStatusUndefined
		}
		return nil, err
	}

	return &models.Status{Id: id, Title: title, Description: description}, nil
}

func (s *Storage) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	var title, description, creatorId string
	var status *models.Status
	var statusId sql.NullInt64
	err := s.db.QueryRowContext(ctx, "SELECT title, description, creatorId, statusId FROM tasks WHERE id = $1", id).Scan(&title, &description, &creatorId, &statusId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, appErrors.ErrTaskNotExists
		}
		return nil, err
	}

	if statusId.Valid {
		status, err = s.GetStatusById(ctx, int(statusId.Int64))
		if err != nil {
			return nil, err
		}
	}

	return &models.Task{Title: title, CreatorId: creatorId, Status: status}, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, title, description string, statusId int) (*models.Status, error) {
	op := "storage.updateStatus"
	log := s.log.With(op)
	var fields []string
	var values []interface{}
	key := 1

	if title != "" {
		fields = append(fields, fmt.Sprintf("title = $%d", key))
		values = append(values, title)
		key++
	}
	if description != "" {
		fields = append(fields, fmt.Sprintf("description = $%d", key))
		values = append(values, description)
		key++
	}

	values = append(values, statusId)

	query := fmt.Sprintf("UPDATE statuses SET %s WHERE id = $%d RETURNING title, description", strings.Join(fields, ", "), key)
	err := s.db.QueryRowContext(ctx, query, values...).Scan(&title, &description)
	if err != nil {
		log.Error("Error: ", err)
		return nil, err
	}

	return &models.Status{Id: statusId, Title: title, Description: description}, nil
}

//TODO: ADD ASSIGNEES
//func (s *Storage) GetAssignedTasksByUID(ctx context.Context, userId int) ([]*models.Task, error) {
//	op := "storage.GetAssignedTasksByUID"
//	log := s.log.With(op)
//	var tasks []*models.Task
//
//	//query the tasks where user is part of assignees
//	rows, err := s.db.QueryContext(ctx, `
//    SELECT t.title, t.description, t.creatorId,
//          t.statusId, t.due, t.completed, ta.userId as assignee_uid,
//          ta.role
//	FROM tasks t
//	JOIN task_assignees ta ON ta.taskId = t.id
//	WHERE t.id IN (
//	SELECT taskId
//	FROM task_assignees
//	WHERE userId = $1)`, userId)
//
//	//close rows on end
//	defer rows.Close()
//
//	if err != nil {
//		log.Error("Error: ", err)
//		return nil, err
//	}
//
//	//get every task by
//	for rows.Next() {
//		var id int
//		var title, creatorId, description string
//		var due time.Time
//		var completed *sql.NullBool
//		var completedWrapper *wrapperspb.BoolValue
//		var status *models.Status
//		var statusId *sql.NullInt64
//
//		// get values from row
//		err = rows.Scan(&title, &description, &creatorId, &statusId, &due, &completed)
//		if err != nil {
//			return nil, err
//		}
//
//		// get status if it is not null
//		if statusId.Valid {
//			status, err = s.GetStatusById(ctx, int(statusId.Int64))
//			if err != nil {
//				return nil, err
//			}
//		}
//
//		// if completed != null
//		if completed.Valid {
//			completedWrapper = wrapperspb.Bool(completed.Bool)
//		}
//
//		if err != nil {
//			return nil, err
//		}
//
//		// generate task model
//		tasks = append(tasks, &models.Task{
//			Id:            id,
//			Title:         title,
//			Description:   description,
//			Due:           due,
//			Completed:     completedWrapper,
//			Status:        status,
//			CreatorId:     creatorId,
//		})
//
//	}
//
//	return tasks, nil
//}

func (s *Storage) GetCreatedTasksByUID(ctx context.Context, userId int) ([]*models.Task, error) {
	op := "storage.GetCreatedTasksByUID"
	log := s.log.With(op)
	var tasks []*models.Task

	//query the tasks where user is part of assignees
	rows, err := s.db.QueryContext(ctx, `
   SELECT t.title, t.description, t.creatorId,
          t.statusId, t.due, t.completed,
		  ta.userId as assignee_uid,
          ta.role
	FROM tasks t
	FULL JOIN task_assignees ta ON ta.taskId = t.id 
	WHERE t.creatorId = $1
	`, userId)

	//close rows on end
	defer rows.Close()

	if err != nil {
		log.Error("Error: ", err)
		return nil, err
	}

	//get every task by
	for rows.Next() {
		var id int
		var title, creatorId, description string
		var due time.Time
		var completed *sql.NullBool
		var completedWrapper *wrapperspb.BoolValue
		var status *models.Status
		var statusId *sql.NullInt64
		//var assignee_id *sql.NullInt64

		// get values from row
		err = rows.Scan(&title, &description, &creatorId, &statusId, &due, &completed)
		if err != nil {
			return nil, err
		}

		// get status if it is not null
		if statusId.Valid {
			status, err = s.GetStatusById(ctx, int(statusId.Int64))
			if err != nil {
				return nil, err
			}
		}

		// if completed != null
		if completed.Valid {
			completedWrapper = wrapperspb.Bool(completed.Bool)
		}

		// generate task model
		tasks = append(tasks, &models.Task{
			Id:          id,
			Title:       title,
			Description: description,
			Due:         due,
			Completed:   completedWrapper,
			Status:      status,
			CreatorId:   creatorId,
		})

	}

	return tasks, nil
}
