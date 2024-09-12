package departmentimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/services/department"

	"go.uber.org/zap"
)

type service struct {
	store *store
	cfg   *config.Config
	log   *zap.Logger
	db    db.DB
}

func NewService(db db.DB, cfg *config.Config) *service {
	return &service{
		store: NewStore(db),
		cfg:   cfg,
		db:    db,
		log:   zap.L().Named("department.service"),
	}
}

func (s *service) CreateDepartment(ctx context.Context, cmd *department.CreateDepartmentCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.departmentTaken(ctx, 0, cmd.Name)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return department.ErrDepartmentAlreadyExists
		}

		err = s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateDepartment(ctx context.Context, cmd *department.UpdateDepartmentCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.departmentTaken(ctx, cmd.ID, cmd.Name)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			return department.ErrDepartmentNotFound
		}

		if len(result) > 1 || (len(result) == 1 && result[0].ID != cmd.ID) {
			return department.ErrDepartmentAlreadyExists
		}

		err = s.store.update(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetDepartmentByID(ctx context.Context, id int) (*department.Department, error) {
	result, err := s.store.getDepartmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) DeleteDepartment(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.getDepartmentByID(ctx, id)
		if err != nil {
			return err
		}

		if result == nil {
			return department.ErrDepartmentNotFound
		}

		err = s.store.delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) SearchDepartment(ctx context.Context, query *department.SearchDepartmentQuery) (*department.SearchDepartmentResult, error) {
	if query.Page <= 0 {
		query.Page = s.cfg.Pagination.Page
	}

	if query.PerPage <= 0 {
		query.PerPage = s.cfg.Pagination.PageLimit
	}

	result, err := s.store.search(ctx, query)
	if err != nil {
		return nil, err
	}

	result.PerPage = query.PerPage
	result.Page = query.Page

	return result, nil
}
