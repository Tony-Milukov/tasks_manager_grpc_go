package models

import "time"

type Task struct {
	Id          int
	Title       string
	Description string
	Due         time.Time
	StatusId    int
	CreatorId   string
	Status      *Status
}

type Status struct {
	Id          int
	Title       string
	Description string
}
