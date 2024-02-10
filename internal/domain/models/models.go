package models

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sso_3.0/internal/domain/user"
	"time"
)

type Task struct {
	Id          int
	Title       string
	Description string
	Due         time.Time
	Completed   *wrapperspb.BoolValue
	CreatorId   string
	Status      *Status
	Assignees   []*Assignee
}

type Status struct {
	Id          int
	Title       string
	Description string
}

type Assignee struct {
	User   *user.Model
	Role   string
	Id     int
	TaskId int
}
type TaskFilters struct {
	AssignedToMe bool
	CreatedByMe  bool
	UnCompleted  bool
	Completed    bool
	AssigneeId   string
	StatusId     int
}
