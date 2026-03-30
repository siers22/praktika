package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type MovementRepository struct {
	pool *pgxpool.Pool
}

func NewMovementRepository(pool *pgxpool.Pool) *MovementRepository {
	return &MovementRepository{pool: pool}
}

const movementSelectBase = `
	SELECT m.id, m.equipment_id, m.from_department_id, m.to_department_id,
	       m.moved_by, m.moved_at, m.reason,
	       e.name, e.inventory_number, fd.name, td.name, u.full_name
	FROM movements m
	LEFT JOIN equipment e ON e.id = m.equipment_id
	LEFT JOIN departments fd ON fd.id = m.from_department_id
	LEFT JOIN departments td ON td.id = m.to_department_id
	LEFT JOIN users u ON u.id = m.moved_by
`

func scanMovement(row interface{ Scan(...any) error }) (*model.Movement, error) {
	m := &model.Movement{}
	return m, row.Scan(
		&m.ID, &m.EquipmentID, &m.FromDepartmentID, &m.ToDepartmentID,
		&m.MovedBy, &m.MovedAt, &m.Reason,
		&m.EquipmentName, &m.InventoryNumber, &m.FromDepartmentName, &m.ToDepartmentName, &m.MovedByName,
	)
}

func (r *MovementRepository) Create(ctx context.Context, mv *model.Movement) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO movements (id, equipment_id, from_department_id, to_department_id, moved_by, moved_at, reason)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, mv.ID, mv.EquipmentID, mv.FromDepartmentID, mv.ToDepartmentID, mv.MovedBy, mv.MovedAt, mv.Reason)
	if err != nil {
		return fmt.Errorf("create movement: %w", err)
	}
	return nil
}

func (r *MovementRepository) ListByEquipment(ctx context.Context, equipmentID uuid.UUID, page, perPage int) ([]*model.Movement, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM movements WHERE equipment_id=$1", equipmentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	rows, err := r.pool.Query(ctx,
		movementSelectBase+" WHERE m.equipment_id=$1 ORDER BY m.moved_at DESC LIMIT $2 OFFSET $3",
		equipmentID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list movements by equipment: %w", err)
	}
	defer rows.Close()

	var list []*model.Movement
	for rows.Next() {
		mv, err := scanMovement(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, mv)
	}
	return list, total, nil
}

func (r *MovementRepository) ListAll(ctx context.Context, dateFrom, dateTo *time.Time, page, perPage int) ([]*model.Movement, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	i := 1
	if dateFrom != nil {
		where += fmt.Sprintf(" AND m.moved_at >= $%d", i)
		args = append(args, dateFrom)
		i++
	}
	if dateTo != nil {
		where += fmt.Sprintf(" AND m.moved_at <= $%d", i)
		args = append(args, dateTo)
		i++
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM movements m "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	args = append(args, perPage, offset)
	rows, err := r.pool.Query(ctx,
		movementSelectBase+where+fmt.Sprintf(" ORDER BY m.moved_at DESC LIMIT $%d OFFSET $%d", i, i+1),
		args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list all movements: %w", err)
	}
	defer rows.Close()

	var list []*model.Movement
	for rows.Next() {
		mv, err := scanMovement(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, mv)
	}
	return list, total, nil
}

func (r *MovementRepository) GetRecent(ctx context.Context, limit int) ([]*model.Movement, error) {
	rows, err := r.pool.Query(ctx,
		movementSelectBase+" ORDER BY m.moved_at DESC LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("get recent movements: %w", err)
	}
	defer rows.Close()
	var list []*model.Movement
	for rows.Next() {
		mv, err := scanMovement(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, mv)
	}
	return list, nil
}
