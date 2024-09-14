package logsmonitoringimpl

import (
	"context"
	"task/config"
	"task/internal/db"
	"task/internal/identity/monitoringactivities/logsmonitoring"

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
		db:    db,
		cfg:   cfg,
		log:   zap.L().Named("logsmonitoring.service"),
	}
}

func (s *service) MonotoringLogs(ctx context.Context, query *logsmonitoring.MonitoringLogsQuery) (*logsmonitoring.MonitoringLogsResult, error) {
	if query.Page <= 0 {
		query.Page = s.cfg.Pagination.Page
	}

	if query.PerPage <= 0 {
		query.PerPage = s.cfg.Pagination.PageLimit
	}

	result, err := s.store.monitoringLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	result.PerPage = query.PerPage
	result.Page = query.Page

	return result, nil
}
