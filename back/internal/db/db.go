package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	slog.Info("connected to database")
	return pool, nil
}

// RunMigrations applies pending *.up.sql migrations from migrationsDir.
func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    BIGINT PRIMARY KEY,
			name       TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	applied := map[int64]bool{}
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return err
		}
		applied[v] = true
	}
	rows.Close()

	pattern := filepath.Join(migrationsDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	sort.Strings(files)

	for _, f := range files {
		base := filepath.Base(f)
		parts := strings.SplitN(base, "_", 2)
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("parse migration version from %s: %w", base, err)
		}

		if applied[version] {
			continue
		}

		sql, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx for migration %s: %w", base, err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("execute migration %s: %w", base, err)
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			version, base,
		); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", base, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", base, err)
		}

		slog.Info("applied migration", "file", base)
	}

	return nil
}
