package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type DepartmentService struct {
	repo      *repository.DepartmentRepository
	auditRepo *repository.AuditRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository, auditRepo *repository.AuditRepository) *DepartmentService {
	return &DepartmentService{repo: repo, auditRepo: auditRepo}
}

func (s *DepartmentService) Create(ctx context.Context, actorID uuid.UUID, req *model.CreateDepartmentRequest) (*model.Department, error) {
	d := &model.Department{
		ID:        uuid.New(),
		Name:      req.Name,
		Location:  req.Location,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	s.auditRepo.Log(ctx, actorID, "create", "department", &d.ID, map[string]any{"name": d.Name})
	return d, nil
}

func (s *DepartmentService) List(ctx context.Context) ([]*model.Department, error) {
	return s.repo.List(ctx)
}

func (s *DepartmentService) Update(ctx context.Context, actorID, id uuid.UUID, req *model.UpdateDepartmentRequest) (*model.Department, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	d.Name = req.Name
	d.Location = req.Location
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	s.auditRepo.Log(ctx, actorID, "update", "department", &id, req)
	return d, nil
}

func (s *DepartmentService) Delete(ctx context.Context, actorID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditRepo.Log(ctx, actorID, "delete", "department", &id, nil)
	return nil
}
