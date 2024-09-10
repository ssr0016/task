package user

import "context"

type Service interface {
	RegisterUser(ctx context.Context, cmd *RegisterUserCommand) error
	CreateUser(ctx context.Context, cmd *CreateUserCommand) error
	GetUserByID(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, cmd *UpdateUserCommand) error
	SearchUser(ctx context.Context, query *SearchUserQuery) (*SearchUserResult, error)
	DeleteUser(ctx context.Context, id int) error
	GetUserByEmail(ctx context.Context, cmd *LoginUserCommand) (string, error)
}
