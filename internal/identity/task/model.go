package task

import (
	"task/internal/api/errors"
)

var (
	ErrTaskAlreadyExists                = errors.New("task.already-exists", "Task already exists")
	ErrTaskNotFound                     = errors.New("task.not-found", "Task not found")
	ErrInvalidTaskTitle                 = errors.New("task.invalid-title", "Invalid task title")
	ErrInvalidTaskDescription           = errors.New("task.invalid-description", "Invalid task description")
	ErrInvalidUserID                    = errors.New("task.invalid-user-id", "Invalid user id")
	ErrInvalidTaskPriority              = errors.New("task.invalid-priority", "Invalid task priority")
	ErrInvalidTaskDifficulty            = errors.New("task.invalid-difficulty", "Invalid task difficulty")
	ErrOnlySuperuserCanApproveTheTask   = errors.New("task.only-superuser-can-approve-the-task", "Only superuser can approve the task")
	TaskIsNotReadyForApproval           = errors.New("task.is-not-ready-for-approval", "Task is not ready for approval")
	ErrOnlyAssignedUserCanSubmitTheTask = errors.New("task.only-assigned-user-can-submit-the-task", "Only assigned user can submit the task")
	TaskIsNotReadyForSubmission         = errors.New("task.is-not-ready-for-submission", "Task is not ready for submission")
	TaskIsNotPending                    = errors.New("task.is-not-pending", "Task is not pending")
)

type TaskStatus int

const (
	TaskPending TaskStatus = iota + 1
	TaskReviewing
	TaskDone
)

var validPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

var validDifficulties = map[string]bool{
	"easy":   true,
	"medium": true,
	"hard":   true,
}

type Task struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status" json:"status"`
	Priority    string     `db:"priority" json:"priority"`
	Difficulty  string     `db:"difficulty" json:"difficulty"`
	UserID      int        `db:"user_id" json:"user_id"`
	CreatedAt   string     `db:"created_at" json:"created_at"`
	UpdatedAt   string     `db:"updated_at" json:"updated_at"`
}

type CreateTaskCommand struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Priority    string     `json:"priority"`
	Difficulty  string     `json:"difficulty"`
	UserID      int        `json:"user_id"`
}

type UpdateTaskCommand struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Difficulty  string     `json:"difficulty"`
	Status      TaskStatus `json:"status"`
	UserID      int        `json:"user_id"`
}

type SearchTaskQuery struct {
	Title       string `query:"title"`
	Description string `query:"description"`
	Status      string `query:"status"`
	Priority    string `query:"priority"`
	Difficulty  string `query:"difficulty"`
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

type SubmitTaskCommand struct {
	TaskID int `json:"task_id"`
	UserID int `json:"user_id"`
}

type ApproveTaskCommand struct {
	TaskID int `json:"task_id"`
	UserID int `json:"user_id"`
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

	if !validPriorities[cmd.Priority] {
		return ErrInvalidTaskPriority
	}

	if !validDifficulties[cmd.Difficulty] {
		return ErrInvalidTaskDifficulty
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

	if !validPriorities[cmd.Priority] {
		return ErrInvalidTaskPriority
	}

	if !validDifficulties[cmd.Difficulty] {
		return ErrInvalidTaskDifficulty
	}

	return nil
}
