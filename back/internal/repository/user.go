package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, full_name, email, role, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, u.ID, u.Username, u.PasswordHash, u.FullName, u.Email, u.Role, u.IsActive, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password_hash, full_name, email, role, is_active, created_at, updated_at
		FROM users WHERE id=$1
	`, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password_hash, full_name, email, role, is_active, created_at, updated_at
		FROM users WHERE username=$1
	`, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context, page, perPage int) ([]*model.User, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.pool.Query(ctx, `
		SELECT id, username, password_hash, full_name, email, role, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET full_name=$1, email=$2, role=$3, updated_at=$4 WHERE id=$5
	`, u.FullName, u.Email, u.Role, time.Now(), u.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE users SET password_hash=$1, updated_at=$2 WHERE id=$3",
		hash, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE users SET is_active=$1, updated_at=$2 WHERE id=$3",
		isActive, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	return nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1,$2,$3,$4)
	`, uuid.New(), userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, error) {
	var userID uuid.UUID
	var expiresAt time.Time
	err := r.pool.QueryRow(ctx,
		"SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash=$1",
		tokenHash,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("get refresh token: %w", err)
	}
	return userID, expiresAt, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash=$1", tokenHash)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *UserRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE expires_at < NOW()")
	return err
}
