package monitoringactivitiesimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/identity/monitoringactivities"

	"go.uber.org/zap"
)

type service struct {
	store *store
	db    db.DB
	cfg   *config.Config
	log   *zap.Logger
}

func NewService(db db.DB, cfg *config.Config) *service {
	return &service{
		store: NewStore(db),
		db:    db,
		cfg:   cfg,
		log:   zap.L().Named("monitoringactivities.service"),
	}
}

func (s *service) LogActivity(ctx context.Context, cmd *monitoringactivities.CreateActivityLogCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		err := s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) SearchLogActivities(ctx context.Context, query *monitoringactivities.SearchLogActivityQuery) (*monitoringactivities.SearchLogActivityResult, error) {
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
