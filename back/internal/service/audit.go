package service

import (
	"context"

	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) List(ctx context.Context, f model.AuditFilter) ([]*model.AuditLog, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PerPage <= 0 || f.PerPage > 50 {
		f.PerPage = 20
	}
	return s.repo.List(ctx, f)
}
