package departmentimpl

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"task/internal/db"
	"task/internal/services/department"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("department.store"),
	}
}

func (s *store) create(ctx context.Context, cmd *department.CreateDepartmentCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			INSERT INTO departments (
				name,
				location
			)VALUES (
				$1,
				$2
			) RETURNING id
		`

		var id int

		err := tx.QueryRow(
			ctx,
			rawSQL,
			cmd.Name,
			cmd.Location,
		).Scan(&id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) getDepartmentByID(ctx context.Context, id int) (*department.Department, error) {
	var department department.Department

	rawSQL := `
		SELECT
			id,
			name,
			location,
			created_at,
			updated_at
		FROM
			departments
		WHERE
			id = $1
	`

	err := s.db.Get(ctx, &department, rawSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, nil
		}
	}

	return &department, nil
}

func (s *store) update(ctx context.Context, cmd *department.UpdateDepartmentCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			UPDATE departments
			SET
				name = $1,
				location = $2
			WHERE
				id = $3
		`

		_, err := tx.Exec(
			ctx,
			rawSQL,
			cmd.Name,
			cmd.Location,
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
			DELETE FROM departments
			WHERE id = $1	
		`

		_, err := tx.Exec(ctx, rawSQL, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) departmentTaken(ctx context.Context, id int, name string) ([]*department.Department, error) {
	var result []*department.Department

	rawSQL := `
		SELECT 
			id,
			name,
			location
		FROM
			departments
		WHERE
			id = $1 OR
			name = $2
	`

	err := s.db.Select(ctx, &result, rawSQL, id, name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) search(ctx context.Context, query *department.SearchDepartmentQuery) (*department.SearchDepartmentResult, error) {
	var (
		result = &department.SearchDepartmentResult{
			Department: make([]*department.Department, 0),
		}
		sql            bytes.Buffer
		whereCondition = make([]string, 0)
		whereParams    = make([]interface{}, 0)
		paramIndex     = 1
	)

	sql.WriteString(`
		SELECT
			id,
			name,
			location,
			created_at,
			updated_at
		FROM
			departments
	`)

	if len(query.Name) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("name ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Name+"%")
		paramIndex++
	}

	if len(query.Location) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("location ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Location+"%")
		paramIndex++
	}

	if len(whereCondition) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereCondition, " AND "))
	}

	sql.WriteString(" ORDER BY id DESC")

	if query.PerPage > 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
		whereParams = append(whereParams, query.PerPage, offset)
	}

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	err = s.db.Select(ctx, &result.Department, sql.String(), whereParams...)
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
