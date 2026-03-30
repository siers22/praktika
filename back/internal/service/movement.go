package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type MovementService struct {
	repo          *repository.MovementRepository
	equipmentRepo *repository.EquipmentRepository
	auditRepo     *repository.AuditRepository
}

func NewMovementService(repo *repository.MovementRepository, equipmentRepo *repository.EquipmentRepository, auditRepo *repository.AuditRepository) *MovementService {
	return &MovementService{repo: repo, equipmentRepo: equipmentRepo, auditRepo: auditRepo}
}

func (s *MovementService) Create(ctx context.Context, actorID uuid.UUID, req *model.CreateMovementRequest) (*model.Movement, error) {
	equipment, err := s.equipmentRepo.GetByID(ctx, req.EquipmentID)
	if err != nil {
		return nil, errors.New("equipment not found")
	}
	if equipment.IsArchived {
		return nil, errors.New("cannot move archived equipment")
	}
	if equipment.DepartmentID == req.ToDepartmentID {
		return nil, errors.New("equipment is already in this department")
	}

	mv := &model.Movement{
		ID:               uuid.New(),
		EquipmentID:      req.EquipmentID,
		FromDepartmentID: equipment.DepartmentID,
		ToDepartmentID:   req.ToDepartmentID,
		MovedBy:          actorID,
		MovedAt:          time.Now(),
		Reason:           req.Reason,
	}

	if err := s.repo.Create(ctx, mv); err != nil {
		return nil, err
	}

	if err := s.equipmentRepo.UpdateDepartment(ctx, req.EquipmentID, req.ToDepartmentID); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, actorID, "move", "equipment", &req.EquipmentID, map[string]any{
		"from": equipment.DepartmentID,
		"to":   req.ToDepartmentID,
	})

	return mv, nil
}

func (s *MovementService) ListByEquipment(ctx context.Context, equipmentID uuid.UUID, page, perPage int) ([]*model.Movement, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	return s.repo.ListByEquipment(ctx, equipmentID, page, perPage)
}

func (s *MovementService) ListAll(ctx context.Context, dateFrom, dateTo *time.Time, page, perPage int) ([]*model.Movement, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	return s.repo.ListAll(ctx, dateFrom, dateTo, page, perPage)
}
