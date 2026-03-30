package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type DepartmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{pool: pool}
}

func (r *DepartmentRepository) Create(ctx context.Context, d *model.Department) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO departments (id, name, location, created_at) VALUES ($1,$2,$3,$4)",
		d.ID, d.Name, d.Location, d.CreatedAt)
	if err != nil {
		return fmt.Errorf("create department: %w", err)
	}
	return nil
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Department, error) {
	d := &model.Department{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, location, created_at FROM departments WHERE id=$1", id,
	).Scan(&d.ID, &d.Name, &d.Location, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get department by id: %w", err)
	}
	return d, nil
}

func (r *DepartmentRepository) List(ctx context.Context) ([]*model.Department, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, location, created_at FROM departments ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	defer rows.Close()
	var list []*model.Department
	for rows.Next() {
		d := &model.Department{}
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, d *model.Department) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE departments SET name=$1, location=$2 WHERE id=$3",
		d.Name, d.Location, d.ID)
	if err != nil {
		return fmt.Errorf("update department: %w", err)
	}
	return nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM departments WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete department: %w", err)
	}
	return nil
}
