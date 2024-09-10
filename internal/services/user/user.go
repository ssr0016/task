package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, cmd *CreateUserCommand) error
	GetUserByID(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, cmd *UpdateUserCommand) error
	SearchUser(ctx context.Context, query *SearchUserQuery) (*SearchUserResult, error)
	DeleteUser(ctx context.Context, id int) error
}
