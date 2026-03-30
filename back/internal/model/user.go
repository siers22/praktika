package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleInventory Role = "inventory"
	RoleViewer    Role = "viewer"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     Role   `json:"role" validate:"required,oneof=admin inventory viewer"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     Role   `json:"role" validate:"required,oneof=admin inventory viewer"`
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   Role      `json:"role"`
}
