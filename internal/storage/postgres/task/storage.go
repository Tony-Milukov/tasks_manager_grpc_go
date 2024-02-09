package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
	"sso_3.0/internal/domain/models"
	"sso_3.0/internal/domain/user"
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

// CreateTask is creating a new tasm with given params
func (s *Storage) CreateTask(ctx context.Context, title, description, creatorId string, statusId int, due time.Time) (*models.Task, error) {
	var id int
	err := s.db.QueryRowContext(ctx, "INSERT INTO tasks (title, description, statusid, creatorId, due) VALUES ($1, $2, $3, $4, $5) RETURNING id", title, description, statusId, creatorId, due).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return s.GetTaskById(ctx, id)
}

// DeleteTask is deleting task by taskId
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

// UpdateTask is updating task by given params where they are not default value
func (s *Storage) UpdateTask(ctx context.Context, title, description string, due time.Time, status *models.Status, completed *wrapperspb.BoolValue, id int) (*models.Task, error) {
	var fields []string
	var values []interface{}
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
	}

	if status != nil {
		fields = append(fields, fmt.Sprintf("statusId = $%d", key))
		values = append(values, status.Id)
		key++
	}

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $1", strings.Join(fields, ", "))

	//execute the update and get new values
	_, err := s.db.ExecContext(ctx, query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return s.GetTaskById(ctx, id)
}

// CreateStatus is creating status with given params
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

// DeleteStatus deletes status by id
// deletes status also from the tasks
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

// GetStatusById gets status by id
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

// GetTaskById gets task by id
func (s *Storage) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	var assignees []*models.Assignee
	op := "storage.GetTaskById"
	log := s.log.With(op)
	var title, creatorId, description string
	var due time.Time
	//completed
	var completed *sql.NullBool
	var completedWrapper *wrapperspb.BoolValue
	//status
	var status *models.Status
	//query the tasks where user is part of assignees
	rows, err := s.db.QueryContext(ctx, `
	SELECT t.id, t.title, t.description, t.creatorId,
		   t.statusId, t.due, t.completed, ta.role, 
		   ta.userId as assigneeId, u.email as assigneeEmail,
		   ta.id as taskId
	FROM tasks t 
	LEFT JOIN task_assignees ta ON ta.taskId = t.id 
	LEFT JOIN users u ON ta.userId = u.id
	WHERE t.id = $1
	`, id)
	if err != nil {
		log.Error("Error: ", err)
		return nil, err
	}
	fmt.Println(rows)
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	//close rows on end
	defer rows.Close()

	//get every task by
	for rows.Next() {
		// define all needed vars

		var statusId *sql.NullInt64
		// assignee
		var assigneeUserId *sql.NullString
		var assigneeId *sql.NullInt64
		var assigneeRole *sql.NullString
		var assigneeEmail *sql.NullString

		// get values from row
		err = rows.Scan(&id, &title, &description, &creatorId, &statusId, &due, &completed, &assigneeRole, &assigneeUserId, &assigneeEmail, &assigneeId)

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
		if completed != nil && completed.Valid {
			completedWrapper = wrapperspb.Bool(completed.Bool)
		}

		if assigneeId != nil && assigneeId.Valid {
			assignees = append(assignees, &models.Assignee{Id: int(assigneeId.Int64), TaskId: id, Role: assigneeRole.String, User: &user.Model{Id: assigneeUserId.String, Email: assigneeEmail.String}})
		}
	}

	return &models.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Due:         due,
		Completed:   completedWrapper,
		Status:      status,
		CreatorId:   creatorId,
		Assignees:   assignees,
	}, nil
}

// UpdateStatus updates status by id with given params
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

// GetCreatedTasksByUID gets tasks by userId
func (s *Storage) GetCreatedTasksByUID(ctx context.Context, userId string) ([]*models.Task, error) {
	op := "storage.GetCreatedTasksByUID"
	log := s.log.With(op)
	var tasks []*models.Task
	assignees := make(map[int][]*models.Assignee)

	//query the tasks where user is part of assignees
	rows, err := s.db.QueryContext(ctx, `
   SELECT t.id, t.title, t.description, t.creatorId,
          t.statusId, t.due, t.completed, ta.role, ta.userId as assigneeId, u.email as assigneeEmail, ta.id
	FROM tasks t
	LEFT JOIN task_assignees ta ON ta.taskId = t.id 
	LEFT JOIN users u ON ta.userId = u.id
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
		// define all needed vars

		//task
		var id int
		var title, creatorId, description string
		var due time.Time
		//completed
		var completed *sql.NullBool
		var completedWrapper *wrapperspb.BoolValue
		//status
		var status *models.Status
		var statusId *sql.NullInt64
		// assignee
		var assigneeUserId *sql.NullString
		var assigneeId *sql.NullInt64
		var assigneeRole *sql.NullString
		var assigneeEmail *sql.NullString

		// get values from row
		err = rows.Scan(&id, &title, &description, &creatorId, &statusId, &due, &completed, &assigneeRole, &assigneeUserId, &assigneeEmail, &assigneeId)
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
		if completed != nil && completed.Valid {
			completedWrapper = wrapperspb.Bool(completed.Bool)
		}

		if assigneeId != nil && assigneeId.Valid {
			assignees[id] = append(assignees[id],
				&models.Assignee{
					Id:     int(assigneeId.Int64),
					TaskId: id,
					Role:   assigneeRole.String,
					User: &user.Model{
						Id:    assigneeUserId.String,
						Email: assigneeEmail.String}})
		}

		if !taskContainsId(tasks, id) {
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
	}

	for _, task := range tasks {
		if assignees[task.Id] != nil {
			task.Assignees = assignees[task.Id]
		}
	}

	return tasks, nil
}

func (s *Storage) AssignTask(ctx context.Context, userId, role string, taskId int) (*models.Task, error) {
	op := "storage.AssignTask"
	log := s.log.With(op)
	var id int
	fmt.Printf("TASK: %d  USER: %s ROLE %s", taskId, userId, role)
	err := s.db.QueryRowContext(ctx, "INSERT INTO task_assignees (taskId, role,userId) VALUES ($1, $2, $3) RETURNING id", taskId, role, userId).Scan(&id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return nil, appErrors.TaskAlreadyAssigned
		}
		log.Error("Error: ", err)
		return nil, err
	}
	return s.GetTaskById(ctx, taskId)
}

// taskContainsId this function checks if given slice of tasks contains given id
func taskContainsId(tasks []*models.Task, id int) bool {
	for _, task := range tasks {
		if task.Id == id {
			return true
		}
	}
	return false
}
