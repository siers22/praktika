package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   *uuid.UUID      `json:"entity_id"`
	Details    json.RawMessage `json:"details"`
	CreatedAt  time.Time       `json:"created_at"`
	// Joined
	Username *string `json:"username,omitempty"`
}

type AuditFilter struct {
	UserID     *uuid.UUID
	Action     string
	EntityType string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PerPage    int
}
