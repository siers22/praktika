package service

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type InventoryService struct {
	repo          *repository.InventoryRepository
	equipmentRepo *repository.EquipmentRepository
	auditRepo     *repository.AuditRepository
}

func NewInventoryService(repo *repository.InventoryRepository, equipmentRepo *repository.EquipmentRepository, auditRepo *repository.AuditRepository) *InventoryService {
	return &InventoryService{repo: repo, equipmentRepo: equipmentRepo, auditRepo: auditRepo}
}

func (s *InventoryService) CreateSession(ctx context.Context, actorID uuid.UUID, req *model.CreateInventorySessionRequest) (*model.InventorySession, error) {
	session := &model.InventorySession{
		ID:           uuid.New(),
		DepartmentID: req.DepartmentID,
		Status:       model.InventoryInProgress,
		CreatedBy:    actorID,
		StartedAt:    time.Now(),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	s.auditRepo.Log(ctx, actorID, "create", "inventory_session", &session.ID, map[string]any{"department_id": req.DepartmentID})
	return s.repo.GetSessionByID(ctx, session.ID)
}

func (s *InventoryService) GetSession(ctx context.Context, id uuid.UUID) (*model.InventorySession, error) {
	return s.repo.GetSessionByID(ctx, id)
}

func (s *InventoryService) ListSessions(ctx context.Context, page, perPage int) ([]*model.InventorySession, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	return s.repo.ListSessions(ctx, page, perPage)
}

func (s *InventoryService) CheckItem(ctx context.Context, actorID, sessionID uuid.UUID, req *model.CheckInventoryItemRequest) (*model.InventoryItem, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.Status == model.InventoryCompleted {
		return nil, errors.New("cannot modify completed inventory session")
	}

	equipment, err := s.equipmentRepo.GetByID(ctx, req.EquipmentID)
	if err != nil {
		return nil, errors.New("equipment not found")
	}

	item := &model.InventoryItem{
		ID:             uuid.New(),
		SessionID:      sessionID,
		EquipmentID:    req.EquipmentID,
		ExpectedStatus: string(equipment.Status),
		ActualStatus:   req.ActualStatus,
		Comment:        req.Comment,
		CheckedAt:      time.Now(),
	}

	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}

	return s.repo.GetItemByID(ctx, item.ID)
}

func (s *InventoryService) UpdateItem(ctx context.Context, actorID, sessionID, itemID uuid.UUID, req *model.UpdateInventoryItemRequest) (*model.InventoryItem, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.Status == model.InventoryCompleted {
		return nil, errors.New("cannot modify completed inventory session")
	}

	item, err := s.repo.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item.SessionID != sessionID {
		return nil, errors.New("item does not belong to this session")
	}

	item.ActualStatus = req.ActualStatus
	item.Comment = req.Comment
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	return s.repo.GetItemByID(ctx, itemID)
}

func (s *InventoryService) GetItems(ctx context.Context, sessionID uuid.UUID) ([]*model.InventoryItem, error) {
	return s.repo.ListItems(ctx, sessionID)
}

func (s *InventoryService) CompleteSession(ctx context.Context, actorID, sessionID uuid.UUID) error {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.Status == model.InventoryCompleted {
		return errors.New("session is already completed")
	}

	if err := s.repo.CompleteSession(ctx, sessionID); err != nil {
		return err
	}
	s.auditRepo.Log(ctx, actorID, "complete", "inventory_session", &sessionID, nil)
	return nil
}

func (s *InventoryService) ExportCSV(ctx context.Context, sessionID uuid.UUID, w io.Writer) error {
	items, err := s.repo.ListItems(ctx, sessionID)
	if err != nil {
		return err
	}

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"Инвентарный номер", "Наименование", "Ожидаемый статус", "Фактический статус", "Комментарий", "Время проверки"})
	for _, item := range items {
		invNum := ""
		if item.InventoryNumber != nil {
			invNum = *item.InventoryNumber
		}
		name := ""
		if item.EquipmentName != nil {
			name = *item.EquipmentName
		}
		comment := ""
		if item.Comment != nil {
			comment = *item.Comment
		}
		_ = cw.Write([]string{
			invNum, name, item.ExpectedStatus, string(item.ActualStatus), comment,
			item.CheckedAt.Format("2006-01-02 15:04:05"),
		})
	}
	cw.Flush()
	return cw.Error()
}
