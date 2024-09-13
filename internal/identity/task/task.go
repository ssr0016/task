package task

import "context"

type Service interface {
	CreateTask(ctx context.Context, cmd *CreateTaskCommand) error
	UpdateTask(ctx context.Context, cmd *UpdateTaskCommand) error
	GetTaskByID(ctx context.Context, id int) (*Task, error)
	DeleteTask(ctx context.Context, id int) error
	SearchTask(ctx context.Context, query *SearchTaskQuery) (*SearchTaskResult, error)
}
