package taskimpl

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"task/internal/db"
	"task/internal/identity/task"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("task.store"),
	}
}

func (s *store) create(ctx context.Context, cmd *task.CreateTaskCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			INSERT INTO tasks (
				title,
				description,
				status,
				priority,
				difficulty,
				user_id
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
			cmd.Title,
			cmd.Description,
			cmd.Status,
			cmd.Priority,
			cmd.Difficulty,
			cmd.UserID,
		).Scan(&id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) update(ctx context.Context, cmd *task.UpdateTaskCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			UPDATE tasks
			SET
				title = $1,
				description = $2,
				status = $3,
				priority = $4,
				difficulty = $5,
				user_id = $6
			WHERE id = $7
		`

		_, err := tx.Exec(
			ctx,
			rawSQL,
			cmd.Title,
			cmd.Description,
			cmd.Status,
			cmd.Priority,
			cmd.Difficulty,
			cmd.UserID,
			cmd.ID,
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) delete(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			DELETE FROM tasks
			WHERE id = $1	
		`

		_, err := tx.Exec(ctx, rawSQL, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) getTaskByID(ctx context.Context, id int) (*task.Task, error) {
	var task task.Task

	rawSQL := `
		SELECT
			id,
			title,
			description,
			status,
			priority,
			difficulty,
			user_id,
			created_at,
			updated_at
		FROM
			tasks
		WHERE
			id = $1
	`

	err := s.db.Get(ctx, &task, rawSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, nil
		}
	}

	return &task, nil
}

func (s *store) taskTaken(ctx context.Context, id int, title string) ([]*task.Task, error) {
	var result []*task.Task

	rawSQL := `
		SELECT
			id,
			title,
			description,
			status,
			priority,
			difficulty,
			user_id,
			created_at,
			updated_at
		FROM
			tasks
		WHERE
			id = $1
			OR title = $2
	`

	err := s.db.Select(ctx, &result, rawSQL, id, title)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) search(ctx context.Context, query *task.SearchTaskQuery) (*task.SearchTaskResult, error) {
	var (
		result = &task.SearchTaskResult{
			Tasks: make([]*task.Task, 0),
		}
		sql            bytes.Buffer
		whereCondition = make([]string, 0)
		whereParams    = make([]interface{}, 0)
		paramIndex     = 1
	)

	sql.WriteString(`
		SELECT
			id,
			title,
			description,
			status,
			priority,
			difficulty,
			user_id,
			created_at,
			updated_at
		FROM
			tasks
	`)

	if len(query.Title) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("title ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Title+"%")
		paramIndex++
	}

	if len(query.Description) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("description ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Description+"%")
		paramIndex++
	}

	if len(query.Status) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("status = $%d", paramIndex))
		whereParams = append(whereParams, query.Status)
		paramIndex++
	}

	if len(query.Priority) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("priority = $%d", paramIndex))
		whereParams = append(whereParams, query.Priority)
		paramIndex++
	}

	if len(query.Difficulty) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("difficulty = $%d", paramIndex))
		whereParams = append(whereParams, query.Difficulty)
		paramIndex++
	}

	if len(whereCondition) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereCondition, " AND "))
	}

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	sql.WriteString(" ORDER BY created_at DESC")

	if query.PerPage > 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
		whereParams = append(whereParams, query.PerPage, offset)
	}

	err = s.db.Select(ctx, &result.Tasks, sql.String(), whereParams...)
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

func (s *store) updateStatus(ctx context.Context, taskID, status int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			UPDATE
				tasks
			SET
				status = $1,
				updated_at = now()
			WHERE
				id = $2
		`

		_, err := tx.Exec(ctx, rawSQL, status, taskID)
		return err
	})
}

func (s *store) isSuperuserOrDefaultUser(ctx context.Context, userID int) (bool, error) {
	var role string

	rawSQL := `
		SELECT
			role
		FROM
			users
		WHERE
			id = $1
	`
	err := s.db.Get(ctx, &role, rawSQL, userID)
	if err != nil {
		return false, err
	}

	return role == "superuser" || role == "user", nil
}
