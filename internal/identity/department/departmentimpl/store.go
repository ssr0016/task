package departmentimpl

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task/internal/db"
	"task/internal/identity/department"
	"task/internal/identity/user"

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

// Assign user to specific department
func (s *store) assignUserToDepartment(ctx context.Context, cmd *department.AssignUserToDepartmentCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			UPDATE users
			SET department_id = $1
			WHERE id = $2
		`
		_, err := tx.Exec(ctx, rawSQL, cmd.DepartmentID, cmd.UserID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) getUsersByDepartment(ctx context.Context, departmentID int) ([]*user.UserDepartmentDTO, error) {
	var users []*user.UserDepartmentDTO

	rawSQL := `
		SELECT
			u.id,
			u.uuid,
			u.first_name,
			u.last_name,
			u.email,
			u.password_hash,
			u.address,
			u.phone_number,
			u.date_of_birth,
			u.role,
			u.status,
			u.created_at,
			u.updated_at,
			u.department_id,
			d.name AS department_name
		FROM
			users u
		LEFT JOIN
			departments d
		ON u.department_id = d.id
		WHERE
			u.department_id = $1
	`

	err := s.db.Select(ctx, &users, rawSQL, departmentID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// RemoveUserFromDepartment removes the department association from the user by setting department_id to NULL
func (s *store) removeUserFromDepartment(ctx context.Context, userID int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			UPDATE users
			SET department_id = NULL
			WHERE id = $1	
		`

		_, err := tx.Exec(ctx, rawSQL, userID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) searchAllUsersByDepartment(ctx context.Context, query *department.SearchAllUsersByDepartmentQuery) (*department.SearchAllUsersByDepartmentResult, error) {
	var (
		result = &department.SearchAllUsersByDepartmentResult{
			User: make([]*user.UserDepartmentDTO, 0),
		}
		sql            bytes.Buffer
		whereCondition = make([]string, 0)
		whereParams    = make([]interface{}, 0)
		paramIndex     = 1
	)

	sql.WriteString(`
		SELECT
			u.id,
			u.uuid,
			u.first_name,
			u.last_name,
			u.email,
			u.password_hash,
			u.address,
			u.phone_number,
			u.date_of_birth,
			u.role,
			u.status,
			u.created_at,
			u.updated_at,
			u.department_id,
			d.name AS department_name
		FROM
			users u
		LEFT JOIN
			departments d
		ON u.department_id = d.id
		WHERE 1=1
	`)

	if len(query.DepartmentName) > 0 {
		whereCondition = append(whereCondition, "d.name ILIKE $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, "%"+query.DepartmentName+"%")
		paramIndex++
	}

	if len(query.Role) > 0 {
		whereCondition = append(whereCondition, "u.role = $"+strconv.Itoa(paramIndex))
		whereParams = append(whereParams, query.Role)
		paramIndex++
	}

	if len(whereCondition) > 0 {
		sql.WriteString(" AND " + strings.Join(whereCondition, " AND "))
	}

	sql.WriteString(" ORDER BY u.id DESC")

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	if query.PerPage > 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
		whereParams = append(whereParams, query.PerPage, offset)
	}

	err = s.db.Select(ctx, &result.User, sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}

	result.TotalCount = count

	return result, nil
}
