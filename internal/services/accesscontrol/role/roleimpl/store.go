package roleimpl

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"task/internal/db"
	"task/internal/services/accesscontrol/role"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("role.store"),
	}
}

func (s *store) create(ctx context.Context, cmd *role.CreateRoleCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		names := strings.Join(cmd.Name, ", ")

		rawSQL := `
			INSERT INTO roles (
				name,
				description	
			)VALUES(
				$1,
				$2
			) RETURNING id
		`
		var id int

		err := tx.QueryRow(ctx, rawSQL, names, cmd.Description).Scan(&id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) getByID(ctx context.Context, id int) (*role.RoleDTO, error) {
	var result role.RoleDTO

	rawSQL := `
        SELECT
            id,
            name,
            description,
			created_at,
			updated_at
        FROM roles
        WHERE id = $1
    `

	err := s.db.Get(ctx, &result, rawSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) getRoles(ctx context.Context) ([]*role.RoleDTO, error) {
	var result []*role.RoleDTO

	rawSQL := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM roles
	`

	err := s.db.Select(ctx, &result, rawSQL)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) delete(ctx context.Context, id int) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			DELETE FROM roles
			WHERE id = $1	
		`

		_, err := tx.Exec(ctx, rawSQL, id)
		if err != nil {
			return err
		}

		return nil
	})

}
