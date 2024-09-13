package taskimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/identity/task"

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
		log:   zap.L().Named("task.service"),
	}
}

func (s *service) CreateTask(ctx context.Context, cmd *task.CreateTaskCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.taskTaken(ctx, 0, cmd.Title)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return task.ErrTaskAlreadyExists
		}

		err = s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateTask(ctx context.Context, cmd *task.UpdateTaskCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.taskTaken(ctx, cmd.ID, cmd.Title)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			return task.ErrTaskNotFound
		}

		if len(result) > 1 || (len(result) == 1 && result[0].ID != cmd.ID) {
			return task.ErrTaskAlreadyExists
		}

		err = s.store.update(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetTaskByID(ctx context.Context, id int) (*task.Task, error) {
	result, err := s.store.getTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, task.ErrTaskNotFound
	}

	return result, nil
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.getTaskByID(ctx, id)
		if err != nil {
			return err
		}

		if result == nil {
			return task.ErrTaskNotFound
		}

		err = s.store.delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) SearchTask(ctx context.Context, query *task.SearchTaskQuery) (*task.SearchTaskResult, error) {
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
