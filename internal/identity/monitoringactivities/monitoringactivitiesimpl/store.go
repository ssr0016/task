package monitoringactivitiesimpl

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"task/internal/db"
	"task/internal/identity/monitoringactivities"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("monitoringactivities.store"),
	}
}

// create logs the activity into the activity_logs table
func (s *store) create(ctx context.Context, cmd *monitoringactivities.CreateActivityLogCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			INSERT INTO activity_logs (
				user_id,
				activity,
				action,
				resource,
				details,
				created_at
			)VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6
			) RETURNING id
		`

		var id int

		err := tx.QueryRow(
			ctx,
			rawSQL,
			cmd.UserID,
			cmd.Activity,
			cmd.Action,
			cmd.Resource,
			cmd.Details,
			cmd.CreatedAt,
		).Scan(&id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) search(ctx context.Context, query *monitoringactivities.SearchLogActivityQuery) (*monitoringactivities.SearchLogActivityResult, error) {
	var (
		result = &monitoringactivities.SearchLogActivityResult{
			Activities: make([]*monitoringactivities.ActivityLog, 0),
		}
		sql            bytes.Buffer
		whereCondition = make([]string, 0)
		whereParams    = make([]interface{}, 0)
		paramIndex     = 1
	)

	sql.WriteString(`
		SELECT
			id,
			user_id,
			activity,
			action,
			resource,
			details,
			created_at
		FROM
			activity_logs
	`)

	if len(query.UserID) > 0 {
		whereCondition = append(whereCondition, "user_id = $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.UserID)
		paramIndex++
	}

	if len(query.Activity) > 0 {
		whereCondition = append(whereCondition, "activity = $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.Activity)
		paramIndex++
	}

	if len(whereCondition) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereCondition, " AND "))
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

	err = s.db.Select(ctx, &result.Activities, sql.String(), whereParams...)
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
