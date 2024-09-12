package department

import (
	"context"
	"task/internal/services/user"
)

type Service interface {
	CreateDepartment(ctx context.Context, cmd *CreateDepartmentCommand) error
	UpdateDepartment(ctx context.Context, cmd *UpdateDepartmentCommand) error
	GetDepartmentByID(ctx context.Context, id int) (*Department, error)
	SearchDepartment(ctx context.Context, query *SearchDepartmentQuery) (*SearchDepartmentResult, error)
	DeleteDepartment(ctx context.Context, id int) error

	// Assign user to specific department
	AssignUserToDepartment(ctx context.Context, cmd *AssignUserToDepartmentCommand) error
	GetUsersByDepartment(ctx context.Context, departmentID int) ([]*user.UserDepartmentDTO, error)
	RemoveUserFromDepartment(ctx context.Context, userID int) error
	SearchAllUsersByDepartment(ctx context.Context, query *SearchAllUsersByDepartmentQuery) (*SearchAllUsersByDepartmentResult, error)
}
