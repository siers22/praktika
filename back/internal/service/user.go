package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repository.UserRepository
	auditRepo *repository.AuditRepository
}

func NewUserService(userRepo *repository.UserRepository, auditRepo *repository.AuditRepository) *UserService {
	return &UserService{userRepo: userRepo, auditRepo: auditRepo}
}

func (s *UserService) Create(ctx context.Context, actorID uuid.UUID, req *model.CreateUserRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &model.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Email:        req.Email,
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, actorID, "create", "user", &user.ID, map[string]any{"username": user.Username})
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, perPage int) ([]*model.User, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	return s.userRepo.List(ctx, page, perPage)
}

func (s *UserService) Update(ctx context.Context, actorID, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Role = req.Role

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, actorID, "update", "user", &id, req)
	return user, nil
}

func (s *UserService) UpdateStatus(ctx context.Context, actorID, id uuid.UUID, isActive bool) error {
	if err := s.userRepo.UpdateStatus(ctx, id, isActive); err != nil {
		return err
	}
	s.auditRepo.Log(ctx, actorID, "update_status", "user", &id, map[string]any{"is_active": isActive})
	return nil
}
