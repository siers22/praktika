package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
}
