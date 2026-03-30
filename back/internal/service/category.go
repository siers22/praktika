package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type CategoryService struct {
	repo      *repository.CategoryRepository
	auditRepo *repository.AuditRepository
}

func NewCategoryService(repo *repository.CategoryRepository, auditRepo *repository.AuditRepository) *CategoryService {
	return &CategoryService{repo: repo, auditRepo: auditRepo}
}

func (s *CategoryService) Create(ctx context.Context, actorID uuid.UUID, req *model.CreateCategoryRequest) (*model.Category, error) {
	c := &model.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.auditRepo.Log(ctx, actorID, "create", "category", &c.ID, map[string]any{"name": c.Name})
	return c, nil
}

func (s *CategoryService) List(ctx context.Context) ([]*model.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Update(ctx context.Context, actorID, id uuid.UUID, req *model.UpdateCategoryRequest) (*model.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Name = req.Name
	c.Description = req.Description
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	s.auditRepo.Log(ctx, actorID, "update", "category", &id, req)
	return c, nil
}

func (s *CategoryService) Delete(ctx context.Context, actorID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditRepo.Log(ctx, actorID, "delete", "category", &id, nil)
	return nil
}
