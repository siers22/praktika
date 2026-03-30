package model

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentStatus string

const (
	StatusInUse      EquipmentStatus = "in_use"
	StatusInStorage  EquipmentStatus = "in_storage"
	StatusInRepair   EquipmentStatus = "in_repair"
	StatusWrittenOff EquipmentStatus = "written_off"
	StatusReserved   EquipmentStatus = "reserved"
)

type Equipment struct {
	ID                    uuid.UUID       `json:"id"`
	InventoryNumber       string          `json:"inventory_number"`
	Name                  string          `json:"name"`
	Description           *string         `json:"description"`
	CategoryID            uuid.UUID       `json:"category_id"`
	SerialNumber          *string         `json:"serial_number"`
	Model                 *string         `json:"model"`
	Manufacturer          *string         `json:"manufacturer"`
	PurchaseDate          *time.Time      `json:"purchase_date"`
	PurchasePrice         *float64        `json:"purchase_price"`
	WarrantyExpiry        *time.Time      `json:"warranty_expiry"`
	Status                EquipmentStatus `json:"status"`
	DepartmentID          uuid.UUID       `json:"department_id"`
	ResponsiblePersonID   *uuid.UUID      `json:"responsible_person_id"`
	Notes                 *string         `json:"notes"`
	IsArchived            bool            `json:"is_archived"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	// Joined
	CategoryName          *string         `json:"category_name,omitempty"`
	DepartmentName        *string         `json:"department_name,omitempty"`
	ResponsiblePersonName *string         `json:"responsible_person_name,omitempty"`
	Photos                []EquipmentPhoto `json:"photos,omitempty"`
}

type EquipmentPhoto struct {
	ID          uuid.UUID `json:"id"`
	EquipmentID uuid.UUID `json:"equipment_id"`
	FilePath    string    `json:"file_path"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type CreateEquipmentRequest struct {
	InventoryNumber     string          `json:"inventory_number" validate:"required,max=50"`
	Name                string          `json:"name" validate:"required,max=255"`
	Description         *string         `json:"description"`
	CategoryID          uuid.UUID       `json:"category_id" validate:"required"`
	SerialNumber        *string         `json:"serial_number"`
	Model               *string         `json:"model"`
	Manufacturer        *string         `json:"manufacturer"`
	PurchaseDate        *time.Time      `json:"purchase_date"`
	PurchasePrice       *float64        `json:"purchase_price"`
	WarrantyExpiry      *time.Time      `json:"warranty_expiry"`
	Status              EquipmentStatus `json:"status" validate:"required,oneof=in_use in_storage in_repair written_off reserved"`
	DepartmentID        uuid.UUID       `json:"department_id" validate:"required"`
	ResponsiblePersonID *uuid.UUID      `json:"responsible_person_id"`
	Notes               *string         `json:"notes"`
}

type UpdateEquipmentRequest struct {
	Name                string          `json:"name" validate:"required,max=255"`
	Description         *string         `json:"description"`
	CategoryID          uuid.UUID       `json:"category_id" validate:"required"`
	SerialNumber        *string         `json:"serial_number"`
	Model               *string         `json:"model"`
	Manufacturer        *string         `json:"manufacturer"`
	PurchaseDate        *time.Time      `json:"purchase_date"`
	PurchasePrice       *float64        `json:"purchase_price"`
	WarrantyExpiry      *time.Time      `json:"warranty_expiry"`
	Status              EquipmentStatus `json:"status" validate:"required,oneof=in_use in_storage in_repair written_off reserved"`
	DepartmentID        uuid.UUID       `json:"department_id" validate:"required"`
	ResponsiblePersonID *uuid.UUID      `json:"responsible_person_id"`
	Notes               *string         `json:"notes"`
}

type EquipmentFilter struct {
	Search       string
	CategoryID   *uuid.UUID
	DepartmentID *uuid.UUID
	Status       *EquipmentStatus
	Page         int
	PerPage      int
}
