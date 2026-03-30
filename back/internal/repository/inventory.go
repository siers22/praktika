package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type InventoryRepository struct {
	pool *pgxpool.Pool
}

func NewInventoryRepository(pool *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{pool: pool}
}

func (r *InventoryRepository) CreateSession(ctx context.Context, s *model.InventorySession) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO inventory_sessions (id, department_id, status, created_by, started_at)
		VALUES ($1,$2,$3,$4,$5)
	`, s.ID, s.DepartmentID, s.Status, s.CreatedBy, s.StartedAt)
	if err != nil {
		return fmt.Errorf("create inventory session: %w", err)
	}
	return nil
}

func (r *InventoryRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*model.InventorySession, error) {
	s := &model.InventorySession{}
	err := r.pool.QueryRow(ctx, `
		SELECT s.id, s.department_id, s.status, s.created_by, s.started_at, s.finished_at,
		       d.name, u.full_name
		FROM inventory_sessions s
		LEFT JOIN departments d ON d.id = s.department_id
		LEFT JOIN users u ON u.id = s.created_by
		WHERE s.id=$1
	`, id).Scan(&s.ID, &s.DepartmentID, &s.Status, &s.CreatedBy, &s.StartedAt, &s.FinishedAt,
		&s.DepartmentName, &s.CreatedByName)
	if err != nil {
		return nil, fmt.Errorf("get inventory session by id: %w", err)
	}
	return s, nil
}

func (r *InventoryRepository) ListSessions(ctx context.Context, page, perPage int) ([]*model.InventorySession, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM inventory_sessions").Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.department_id, s.status, s.created_by, s.started_at, s.finished_at,
		       d.name, u.full_name
		FROM inventory_sessions s
		LEFT JOIN departments d ON d.id = s.department_id
		LEFT JOIN users u ON u.id = s.created_by
		ORDER BY s.started_at DESC LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list inventory sessions: %w", err)
	}
	defer rows.Close()

	var list []*model.InventorySession
	for rows.Next() {
		s := &model.InventorySession{}
		if err := rows.Scan(&s.ID, &s.DepartmentID, &s.Status, &s.CreatedBy, &s.StartedAt, &s.FinishedAt,
			&s.DepartmentName, &s.CreatedByName); err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, total, nil
}

func (r *InventoryRepository) CompleteSession(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory_sessions SET status='completed', finished_at=$1 WHERE id=$2",
		now, id)
	if err != nil {
		return fmt.Errorf("complete inventory session: %w", err)
	}
	return nil
}

func (r *InventoryRepository) CreateItem(ctx context.Context, item *model.InventoryItem) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO inventory_items (id, session_id, equipment_id, expected_status, actual_status, comment, checked_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, item.ID, item.SessionID, item.EquipmentID, item.ExpectedStatus, item.ActualStatus, item.Comment, item.CheckedAt)
	if err != nil {
		return fmt.Errorf("create inventory item: %w", err)
	}
	return nil
}

func (r *InventoryRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*model.InventoryItem, error) {
	item := &model.InventoryItem{}
	err := r.pool.QueryRow(ctx, `
		SELECT ii.id, ii.session_id, ii.equipment_id, ii.expected_status, ii.actual_status, ii.comment, ii.checked_at,
		       e.name, e.inventory_number
		FROM inventory_items ii
		LEFT JOIN equipment e ON e.id = ii.equipment_id
		WHERE ii.id=$1
	`, id).Scan(&item.ID, &item.SessionID, &item.EquipmentID, &item.ExpectedStatus,
		&item.ActualStatus, &item.Comment, &item.CheckedAt, &item.EquipmentName, &item.InventoryNumber)
	if err != nil {
		return nil, fmt.Errorf("get inventory item by id: %w", err)
	}
	return item, nil
}

func (r *InventoryRepository) ListItems(ctx context.Context, sessionID uuid.UUID) ([]*model.InventoryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ii.id, ii.session_id, ii.equipment_id, ii.expected_status, ii.actual_status, ii.comment, ii.checked_at,
		       e.name, e.inventory_number
		FROM inventory_items ii
		LEFT JOIN equipment e ON e.id = ii.equipment_id
		WHERE ii.session_id=$1
		ORDER BY ii.checked_at
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list inventory items: %w", err)
	}
	defer rows.Close()
	var list []*model.InventoryItem
	for rows.Next() {
		item := &model.InventoryItem{}
		if err := rows.Scan(&item.ID, &item.SessionID, &item.EquipmentID, &item.ExpectedStatus,
			&item.ActualStatus, &item.Comment, &item.CheckedAt, &item.EquipmentName, &item.InventoryNumber); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *InventoryRepository) UpdateItem(ctx context.Context, item *model.InventoryItem) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory_items SET actual_status=$1, comment=$2, checked_at=$3 WHERE id=$4",
		item.ActualStatus, item.Comment, time.Now(), item.ID)
	if err != nil {
		return fmt.Errorf("update inventory item: %w", err)
	}
	return nil
}
