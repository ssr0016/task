package task

import (
	"task/internal/api/errors"
)

var (
	ErrTaskAlreadyExists      = errors.New("task.already-exists", "Task already exists")
	ErrTaskNotFound           = errors.New("task.not-found", "Task not found")
	ErrInvalidTaskTitle       = errors.New("task.invalid-title", "Invalid task title")
	ErrInvalidTaskDescription = errors.New("task.invalid-description", "Invalid task description")
	ErrInvalidUserID          = errors.New("task.invalid-user-id", "Invalid user id")
)

type TaskStatus int

const (
	TaskPending TaskStatus = iota + 1
	TaskInProgress
	TaskDone
)

type Task struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status" json:"status"`
	UserID      int        `db:"user_id" json:"user_id"`
	CreatedAt   string     `db:"created_at" json:"created_at"`
	UpdatedAt   string     `db:"updated_at" json:"updated_at"`
}

type CreateTaskCommand struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      int        `json:"user_id"`
}

type UpdateTaskCommand struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      int        `json:"user_id"`
}

type SearchTaskQuery struct {
	Title       string `query:"title"`
	Description string `query:"description"`
	Status      string `query:"status"`
	UserID      int    `query:"user_id"`
	Page        int    `query:"page"`
	PerPage     int    `query:"per_page"`
}

type SearchTaskResult struct {
	TotalCount int     `json:"total_count"`
	Tasks      []*Task `json:"result"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
}

func (cmd *CreateTaskCommand) Validate() error {
	if len(cmd.Title) == 0 || len(cmd.Title) <= 2 {
		return ErrInvalidTaskTitle
	}

	if len(cmd.Description) == 0 || len(cmd.Description) <= 2 {
		return ErrInvalidTaskDescription
	}

	if cmd.UserID <= 0 {
		return ErrInvalidUserID
	}

	return nil
}

func (cmd *UpdateTaskCommand) Validate() error {
	if cmd.ID <= 0 {
		return ErrInvalidTaskTitle
	}

	if len(cmd.Title) == 0 || len(cmd.Title) <= 2 {
		return ErrInvalidTaskTitle
	}

	if len(cmd.Description) == 0 || len(cmd.Description) <= 2 {
		return ErrInvalidTaskDescription
	}

	if cmd.UserID <= 0 {
		return ErrInvalidUserID
	}

	return nil
}
