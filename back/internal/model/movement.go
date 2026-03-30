package model

import (
	"time"

	"github.com/google/uuid"
)

type Movement struct {
	ID               uuid.UUID `json:"id"`
	EquipmentID      uuid.UUID `json:"equipment_id"`
	FromDepartmentID uuid.UUID `json:"from_department_id"`
	ToDepartmentID   uuid.UUID `json:"to_department_id"`
	MovedBy          uuid.UUID `json:"moved_by"`
	MovedAt          time.Time `json:"moved_at"`
	Reason           *string   `json:"reason"`
	// Joined
	EquipmentName      *string `json:"equipment_name,omitempty"`
	InventoryNumber    *string `json:"inventory_number,omitempty"`
	FromDepartmentName *string `json:"from_department_name,omitempty"`
	ToDepartmentName   *string `json:"to_department_name,omitempty"`
	MovedByName        *string `json:"moved_by_name,omitempty"`
}

type CreateMovementRequest struct {
	EquipmentID    uuid.UUID `json:"equipment_id" validate:"required"`
	ToDepartmentID uuid.UUID `json:"to_department_id" validate:"required"`
	Reason         *string   `json:"reason"`
}

type MovementFilter struct {
	EquipmentID  *uuid.UUID
	DateFrom     *time.Time
	DateTo       *time.Time
	Page         int
	PerPage      int
}
