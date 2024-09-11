package userimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/services/user"
	"task/pkg/util/jwt"
	util "task/pkg/util/password"

	"go.uber.org/zap"
)

type service struct {
	store *store
	cfg   *config.Config
	log   *zap.Logger
	db    db.DB
}

func NewService(db db.DB, cfg *config.Config) service {
	return service{
		store: NewStore(db),
		cfg:   cfg,
		db:    db,
		log:   zap.L().Named("user.service"),
	}
}

func (s *service) CreateUser(ctx context.Context, cmd *user.CreateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.userTaken(ctx, 0, cmd.Email)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return user.ErrUserAlreadyExists
		}

		passwordHash, err := util.HashPassword(cmd.Password)
		if err != nil {
			return err
		}

		cmd.Password = passwordHash

		err = s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateUser(ctx context.Context, cmd *user.UpdateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		// Check if the user exists
		existingUser, err := s.store.getUserByID(ctx, cmd.ID)
		if err != nil {
			return err
		}

		if existingUser == nil {
			return user.ErrUserNotFound
		}

		// Check if the email is already taken
		emailExists, err := s.store.getUserByEmail(ctx, cmd.Email)
		if err != nil {
			return err
		}

		if emailExists != nil && emailExists.ID != cmd.ID {
			return user.ErrUserAlreadyExists
		}

		err = s.store.updateUser(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetUserByID(ctx context.Context, id int) (*user.User, error) {
	result, err := s.store.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) SearchUser(ctx context.Context, query *user.SearchUserQuery) (*user.SearchUserResult, error) {
	if query.Page <= 0 {
		query.Page = s.cfg.Pagination.Page
	}

	if query.PerPage <= 0 {
		query.PerPage = s.cfg.Pagination.PageLimit
	}

	result, err := s.store.searchUser(ctx, query)
	if err != nil {
		return nil, err
	}

	result.PerPage = query.PerPage
	result.Page = query.Page

	return result, nil
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.getUserByID(ctx, id)
		if err != nil {
			return err
		}

		if result == nil {
			return user.ErrUserNotFound
		}

		err = s.store.deleteUser(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetUserByEmail(ctx context.Context, cmd *user.LoginUserCommand) (string, error) {
	result, err := s.store.getUserByEmail(ctx, cmd.Email)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", user.ErrUserNotFound
	}

	// Check if the password is correct
	err = util.CheckPasswordHash(result.PasswordHash, cmd.Password)
	if err != nil {
		return "", user.ErrInvalidPassword
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(result.Email, result.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RegisterUser(ctx context.Context, cmd *user.RegisterUserCommand) error {
	// Ensuring that the user role is set
	role := "user"

	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.userTaken(ctx, 0, cmd.Email)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return user.ErrUserAlreadyExists
		}

		passwordHash, err := util.HashPassword(cmd.Password)
		if err != nil {
			return err
		}

		cmd.Password = passwordHash

		err = s.store.registerUser(ctx, cmd, role)
		if err != nil {
			return err
		}

		return nil
	})

}
