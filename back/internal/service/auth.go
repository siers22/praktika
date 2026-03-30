package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/config"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type jwtClaims struct {
	UserID uuid.UUID  `json:"user_id"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  *repository.UserRepository
	auditRepo *repository.AuditRepository
	cfg       *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, auditRepo *repository.AuditRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, auditRepo: auditRepo, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.TokenPair, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		slog.Debug("login failed: user not found", "username", req.Username)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		slog.Debug("login failed: wrong password", "username", req.Username)
		return nil, ErrInvalidCredentials
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, user.ID, "login", "user", &user.ID, nil)
	slog.Info("user logged in", "user_id", user.ID, "username", user.Username)
	return pair, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	hash := hashToken(refreshToken)
	userID, expiresAt, err := s.userRepo.GetRefreshToken(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(expiresAt) {
		_ = s.userRepo.DeleteRefreshToken(ctx, hash)
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || !user.IsActive {
		return nil, ErrInvalidToken
	}

	_ = s.userRepo.DeleteRefreshToken(ctx, hash)
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(ctx, hash)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &model.Claims{UserID: claims.UserID, Role: claims.Role}, nil
}

func (s *AuthService) issueTokens(ctx context.Context, user *model.User) (*model.TokenPair, error) {
	now := time.Now()

	claims := jwtClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTAccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshRaw, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshHash := hashToken(refreshRaw)
	expiresAt := now.Add(s.cfg.JWTRefreshTTL)

	if err := s.userRepo.SaveRefreshToken(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
	}, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
