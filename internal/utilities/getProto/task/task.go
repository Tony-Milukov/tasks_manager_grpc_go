package task

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"sso_3.0/internal/domain/models"
	protoStatus "sso_3.0/internal/utilities/getProto/status"
	api "sso_3.0/proto/gen"
)

func GetProtoTask(task *models.Task) *api.Task {
	var status *api.Status
	var completed bool
	var assignees []*api.TaskAssignee

	if task == nil {
		return nil
	}

	if task.Status != nil {
		status = protoStatus.GetStatus(task.Status)
	}

	if task.Completed != nil {
		completed = task.Completed.Value
	}

	if task.Assignees != nil {
		assignees = GetProtoAssignees(task.Assignees)
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

func GetProtoTasks(tasks []*models.Task) []*api.Task {
	var value []*api.Task

	for _, task := range tasks {
		protoTask := GetProtoTask(task)
		value = append(value, protoTask)
	}

	return value
}

func GetProtoAssignees(assignees []*models.Assignee) []*api.TaskAssignee {
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
