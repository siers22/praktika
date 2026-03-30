package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdminUser creates a default admin user if no users exist yet.
func SeedAdminUser(pool *pgxpool.Pool, username, password string) error {
	ctx := context.Background()

	var count int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already have users, skip
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, full_name, email, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'admin', true, $6, $7)
	`, uuid.New(), username, string(hash), "Administrator", username+"@localhost", now, now)
	if err != nil {
		return err
	}

	slog.Info("default admin user created", "username", username)
	return nil
}
