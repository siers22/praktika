package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO categories (id, name, description, created_at) VALUES ($1,$2,$3,$4)",
		c.ID, c.Name, c.Description, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	c := &model.Category{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, description, created_at FROM categories WHERE id=$1", id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	return c, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, description, created_at FROM categories ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()
	var list []*model.Category
	for rows.Next() {
		c := &model.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *model.Category) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE categories SET name=$1, description=$2 WHERE id=$3",
		c.Name, c.Description, c.ID)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
