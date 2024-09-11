package role

import (
	"task/internal/api/errors"
)

var (
	ErrNameIsEmpty       = errors.New("role.name-is-empty", "Name is empty")
	ErrNamesIsEmpty      = errors.New("role.names-is-empty", "Names is empty")
	ErrInvalidID         = errors.New("role.invalid-id", "Invalid id")
	ErrRoleAlreadyExists = errors.New("role.already-exists", "Role already exists")
	ErrRoleNotFound      = errors.New("role.not-found", "Role not found")
)

type Role struct {
	ID          int      `db:"id" json:"id"`
	Name        []string `db:"name" json:"name"`
	Description string   `db:"description" json:"description"`
	CreatedAt   string   `db:"created_at" json:"created_at"`
	UpdatedAt   string   `db:"updated_at" json:"updated_at"`
}

type RoleDTO struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type CreateRoleCommand struct {
	Name        []string `json:"name"`
	Description string   `json:"description"`
}

func (cmd *CreateRoleCommand) Validate() error {
	if len(cmd.Name) == 0 {
		return ErrNameIsEmpty
	}

	for _, name := range cmd.Name {
		if len(name) == 0 {
			return ErrNameIsEmpty
		}

		if inValidNames(name) {
			return ErrNamesIsEmpty
		}
	}

	return nil
}

func inValidNames(name string) bool {
	validNames := map[string]bool{
		"user":      true,
		"superuser": true,
	}

	return !validNames[name]
}
