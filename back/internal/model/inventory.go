package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryStatus string

const (
	InventoryInProgress InventoryStatus = "in_progress"
	InventoryCompleted  InventoryStatus = "completed"
)

type ActualStatus string

const (
	ActualFound    ActualStatus = "found"
	ActualNotFound ActualStatus = "not_found"
	ActualDamaged  ActualStatus = "damaged"
)

type InventorySession struct {
	ID             uuid.UUID       `json:"id"`
	DepartmentID   uuid.UUID       `json:"department_id"`
	Status         InventoryStatus `json:"status"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	StartedAt      time.Time       `json:"started_at"`
	FinishedAt     *time.Time      `json:"finished_at"`
	// Joined
	DepartmentName *string         `json:"department_name,omitempty"`
	CreatedByName  *string         `json:"created_by_name,omitempty"`
	ItemsTotal     int             `json:"items_total,omitempty"`
	ItemsChecked   int             `json:"items_checked,omitempty"`
}

type InventoryItem struct {
	ID             uuid.UUID    `json:"id"`
	SessionID      uuid.UUID    `json:"session_id"`
	EquipmentID    uuid.UUID    `json:"equipment_id"`
	ExpectedStatus string       `json:"expected_status"`
	ActualStatus   ActualStatus `json:"actual_status"`
	Comment        *string      `json:"comment"`
	CheckedAt      time.Time    `json:"checked_at"`
	// Joined
	EquipmentName   *string     `json:"equipment_name,omitempty"`
	InventoryNumber *string     `json:"inventory_number,omitempty"`
}

type CreateInventorySessionRequest struct {
	DepartmentID uuid.UUID `json:"department_id" validate:"required"`
}

type CheckInventoryItemRequest struct {
	EquipmentID    uuid.UUID    `json:"equipment_id" validate:"required"`
	ActualStatus   ActualStatus `json:"actual_status" validate:"required,oneof=found not_found damaged"`
	Comment        *string      `json:"comment"`
}

type UpdateInventoryItemRequest struct {
	ActualStatus ActualStatus `json:"actual_status" validate:"required,oneof=found not_found damaged"`
	Comment      *string      `json:"comment"`
}
