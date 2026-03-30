package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, details, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, log.ID, log.UserID, log.Action, log.EntityType, log.EntityID, log.Details, log.CreatedAt)
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) List(ctx context.Context, f model.AuditFilter) ([]*model.AuditLog, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if f.UserID != nil {
		where += fmt.Sprintf(" AND al.user_id=$%d", i)
		args = append(args, *f.UserID)
		i++
	}
	if f.Action != "" {
		where += fmt.Sprintf(" AND al.action=$%d", i)
		args = append(args, f.Action)
		i++
	}
	if f.EntityType != "" {
		where += fmt.Sprintf(" AND al.entity_type=$%d", i)
		args = append(args, f.EntityType)
		i++
	}
	if f.DateFrom != nil {
		where += fmt.Sprintf(" AND al.created_at >= $%d", i)
		args = append(args, f.DateFrom)
		i++
	}
	if f.DateTo != nil {
		where += fmt.Sprintf(" AND al.created_at <= $%d", i)
		args = append(args, f.DateTo)
		i++
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM audit_logs al "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	perPage := f.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage
	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, `
		SELECT al.id, al.user_id, al.action, al.entity_type, al.entity_id, al.details, al.created_at, u.username
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		`+where+fmt.Sprintf(" ORDER BY al.created_at DESC LIMIT $%d OFFSET $%d", i, i+1),
		args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	var list []*model.AuditLog
	for rows.Next() {
		l := &model.AuditLog{}
		var details []byte
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID, &details, &l.CreatedAt, &l.Username); err != nil {
			return nil, 0, err
		}
		if details != nil {
			l.Details = json.RawMessage(details)
		}
		list = append(list, l)
	}
	return list, total, nil
}

func (r *AuditRepository) Log(ctx context.Context, userID uuid.UUID, action, entityType string, entityID *uuid.UUID, details any) {
	var raw json.RawMessage
	if details != nil {
		b, err := json.Marshal(details)
		if err == nil {
			raw = b
		}
	}
	entry := &model.AuditLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    raw,
	}
	// Best-effort: don't fail the main operation on audit error
	_ = r.Create(ctx, entry)
}
