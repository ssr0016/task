package logsmonitoringimpl

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"task/internal/db"
	"task/internal/identity/monitoringactivities/logsmonitoring"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("logsmonitoring.store"),
	}
}

func (s *store) monitoringLogs(ctx context.Context, query *logsmonitoring.MonitoringLogsQuery) (*logsmonitoring.MonitoringLogsResult, error) {
	var (
		result = &logsmonitoring.MonitoringLogsResult{
			Logs: make([]*logsmonitoring.MonitoringLogs, 0),
		}
		sql             bytes.Buffer
		whereConditions = make([]string, 0)
		whereParams     = make([]interface{}, 0)
		paramIndex      = 1
	)

	sql.WriteString(`
		SELECT
			id,
			user_id,
			created_at
		FROM
			activity_logs
	`)

	if len(query.UserID) > 0 {
		whereConditions = append(whereConditions, "user_id = $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.UserID)
		paramIndex++
	}

	if query.DateFrom != nil {
		whereConditions = append(whereConditions, "created_at >= $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.DateFrom)
		paramIndex++
	}

	if query.DateTo != nil {
		whereConditions = append(whereConditions, "created_at <= $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.DateTo)
		paramIndex++
	}

	if len(whereConditions) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereConditions, " AND "))
	}

	sql.WriteString(" ORDER BY created_at DESC")

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	if query.PerPage > 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
		whereParams = append(whereParams, query.PerPage, offset)
	}

	err = s.db.Select(ctx, &result.Logs, sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}

	result.TotalCount = count

	return result, nil
}

func (s *store) getCount(ctx context.Context, sql bytes.Buffer, whereParams []interface{}) (int, error) {
	var count int

	rawSQL := "SELECT COUNT(*) FROM (" + sql.String() + ") as t1"

	err := s.db.Get(ctx, &count, rawSQL, whereParams...)
	if err != nil {
		return 0, err
	}

	return count, nil
}
