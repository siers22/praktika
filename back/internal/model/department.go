package model

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Location  *string   `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDepartmentRequest struct {
	Name     string  `json:"name" validate:"required,max=150"`
	Location *string `json:"location"`
}

type UpdateDepartmentRequest struct {
	Name     string  `json:"name" validate:"required,max=150"`
	Location *string `json:"location"`
}
