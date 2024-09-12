package department

import (
	"task/internal/api/errors"
)

var (
	ErrDepartmentAlreadyExists = errors.New("department.already-exists", "Department already exists")
	ErrDepartmentNotFound      = errors.New("department.not-found", "Department not found")
	ErrInvalidDepartmentName   = errors.New("department.invalid-name", "Invalid department name")
)

type Department struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Location  string `db:"location" json:"location"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type CreateDepartmentCommand struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type UpdateDepartmentCommand struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type SearchDepartmentQuery struct {
	Name     string `query:"name"`
	Location string `query:"location"`
	Page     int    `query:"page"`
	PerPage  int    `query:"per_page"`
}

type SearchDepartmentResult struct {
	TotalCount int           `json:"total_count"`
	Department []*Department `json:"result"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
}

func (cmd *CreateDepartmentCommand) Validate() error {
	if len(cmd.Name) == 0 || len(cmd.Name) <= 2 {
		return ErrInvalidDepartmentName
	}

	return nil
}

func (cmd *UpdateDepartmentCommand) Validate() error {
	if cmd.ID <= 0 {
		return ErrDepartmentNotFound
	}

	if len(cmd.Name) == 0 || len(cmd.Name) <= 2 {
		return ErrInvalidDepartmentName
	}

	return nil
}
