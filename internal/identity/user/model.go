package user

import (
	"strings"
	"task/internal/api/errors"
	util "task/pkg/util/password"
	"task/pkg/util/validation"
	"time"
)

var (
	ErrInvalidEmail       = errors.New("user.invalid-email", "Invalid email")
	ErrInvalidID          = errors.New("user.invalid-id", "Invalid id")
	ErrUserAlreadyExists  = errors.New("user.already-exists", "User already exists")
	ErrUserNotFound       = errors.New("user.not-found", "User not found")
	ErrInvalidPassword    = errors.New("user.invalid-password", "Invalid password")
	ErrInvalidFirstName   = errors.New("user.invalid-first-name", "Invalid first name")
	ErrInvalidLastName    = errors.New("user.invalid-last-name", "Invalid last name")
	ErrInvalidAddress     = errors.New("user.invalid-address", "Invalid address")
	ErrInvalidPhoneNumber = errors.New("user.invalid-phone-number", "Invalid phone number")
	ErrInvalidDateOfBirth = errors.New("user.invalid-date-of-birth", "Invalid date of birth")
	ErrEmailAlreadyExists = errors.New("user.email-already-exists", "Email already exists")
	ErrorInvalidRole      = errors.New("user.invalid-role", "Invalid role")
	ErrInvalidStatus      = errors.New("user.invalid-status", "Invalid status")
)

type Status int

const (
	Active Status = iota + 1
	Inactive
	Deleted
)

const (
	RoleUser      = "user"
	RoleSuperUser = "superuser"
	RoleHR        = "hr"
	RoleManager   = "manager"
)

var validRoles = map[string]bool{
	RoleUser:      true,
	RoleSuperUser: true,
	RoleHR:        true,
	RoleManager:   true,
}

// IsValidRole checks if a role is valid.
func IsValidRole(role string) bool {
	return validRoles[role]
}

type User struct {
	ID           int       `db:"id" json:"id"`
	UUID         string    `db:"uuid" json:"uuid"` // UUID for global uniqueness
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"` // Exclude from JSON output
	Address      string    `db:"address" json:"address"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	DateOfBirth  string    `db:"date_of_birth" json:"date_of_birth"`
	Role         string    `db:"role" json:"role"`
	Status       Status    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"` // Timestamp for creation
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"` // Timestamp for updates
}

type UserDepartmentDTO struct {
	ID           int       `db:"id" json:"id"`
	UUID         string    `db:"uuid" json:"uuid"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Address      string    `db:"address" json:"address"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	DateOfBirth  string    `db:"date_of_birth" json:"date_of_birth"`
	Role         string    `db:"role" json:"role"`
	Status       Status    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	DepartmentID   *int    `db:"department_id" json:"department_id"`
	DepartmentName *string `db:"department_name" json:"department_name"`
}

type CreateUserCommand struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Role        string `json:"role"`
	Status      Status `json:"status"`
}

type UpdateUserCommand struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Role        string `json:"role"`
	Status      Status `json:"status"`
}

type SearchUserQuery struct {
	FirstName   string `query:"first_name"`
	LastName    string `query:"last_name"`
	Email       string `query:"email"`
	Address     string `query:"address"`
	PhoneNumber string `query:"phone_number"`
	DateOfBirth string `query:"date_of_birth"`
	Role        string `query:"role"`
	Status      Status `query:"status"`
	Page        int    `query:"page"`
	PerPage     int    `query:"per_page"`
}

type SearchUserResult struct {
	TotalCount int     `json:"total_count"`
	User       []*User `json:"result"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
}

type LoginUserCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogutUserCommand struct {
	Token string `json:"token"`
}

type RegisterUserCommand struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Status      Status `json:"status"`
}

// Validation for CreateUserCommand
func (cmd *CreateUserCommand) Validate() error {
	if len(cmd.FirstName) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(cmd.LastName) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	if len(cmd.Address) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.DateOfBirth) == 0 {
		return ErrInvalidDateOfBirth
	}
	if !IsValidRole(cmd.Role) {
		return ErrorInvalidRole
	}
	if cmd.Status != Active && cmd.Status != Inactive && cmd.Status != Deleted {
		return ErrInvalidStatus
	}
	return nil
}

// Validation for UpdateUserCommand
func (cmd *UpdateUserCommand) Validate() error {
	if cmd.ID == 0 {
		return ErrUserNotFound
	}
	if len(strings.TrimSpace(cmd.FirstName)) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(strings.TrimSpace(cmd.LastName)) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(strings.TrimSpace(cmd.Address)) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 || !validation.IsValidPhoneNumber(cmd.PhoneNumber) {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.Email) > 0 && !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if !IsValidRole(cmd.Role) {
		return ErrorInvalidRole
	}
	return nil
}

// Validation for LoginUserCommand
func (cmd *LoginUserCommand) Validate() error {
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	return nil
}

// Validation for RegisterUserCommand
func (cmd *RegisterUserCommand) Validate() error {
	if len(cmd.FirstName) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(cmd.LastName) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	if len(cmd.Address) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.DateOfBirth) == 0 {
		return ErrInvalidDateOfBirth
	}
	return nil
}
