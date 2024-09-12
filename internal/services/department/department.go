package department

import "context"

type Service interface {
	CreateDepartment(ctx context.Context, cmd *CreateDepartmentCommand) error
	UpdateDepartment(ctx context.Context, cmd *UpdateDepartmentCommand) error
	GetDepartmentByID(ctx context.Context, id int) (*Department, error)
	SearchDepartment(ctx context.Context, query *SearchDepartmentQuery) (*SearchDepartmentResult, error)
	DeleteDepartment(ctx context.Context, id int) error
}
