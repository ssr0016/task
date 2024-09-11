package roleimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/services/accesscontrol/role"

	"go.uber.org/zap"
)

type service struct {
	store  *store
	cfg    *config.Config
	logger *zap.Logger
	db     db.DB
}

func NewService(db db.DB, cfg *config.Config) service {
	return service{
		store:  NewStore(db),
		db:     db,
		cfg:    cfg,
		logger: zap.L().Named("role.service"),
	}
}

func (s *service) CreateRole(ctx context.Context, cmd *role.CreateRoleCommand) error {
	return s.store.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		err := s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetRoleByID(ctx context.Context, id int) (*role.RoleDTO, error) {
	result, err := s.store.getByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, role.ErrRoleNotFound
	}

	return result, nil
}

func (s *service) GetRoles(ctx context.Context) ([]*role.RoleDTO, error) {
	result, err := s.store.getRoles(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) DeleteRole(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.getByID(ctx, id)
		if err != nil {
			return err
		}

		if result == nil {
			return role.ErrRoleNotFound
		}

		err = s.store.delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}
