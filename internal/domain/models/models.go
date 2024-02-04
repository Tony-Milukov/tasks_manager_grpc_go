package models

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sso_3.0/internal/domain/user"
	"time"
)

type Task struct {
	Id            int
	Title         string
	Description   string
	Due           time.Time
	Completed     *wrapperspb.BoolValue
	CreatorId     string
	Status        *Status
	MajorAssignee string
	MinorAssignee string
}

type Status struct {
	Id          int
	Title       string
	Description string
}

type Assignee struct {
	User *user.Model
	Role string
	Id   int
}

type Assignees struct {
	Major *Assignee
	Minor *Assignee
}
