package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siers22/praktika/back/internal/model"
)

type EquipmentRepository struct {
	pool *pgxpool.Pool
}

func NewEquipmentRepository(pool *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{pool: pool}
}

const equipmentSelectBase = `
	SELECT e.id, e.inventory_number, e.name, e.description, e.category_id,
	       e.serial_number, e.model, e.manufacturer, e.purchase_date, e.purchase_price,
	       e.warranty_expiry, e.status, e.department_id, e.responsible_person_id,
	       e.notes, e.is_archived, e.created_at, e.updated_at,
	       c.name, d.name, u.full_name
	FROM equipment e
	LEFT JOIN categories c ON c.id = e.category_id
	LEFT JOIN departments d ON d.id = e.department_id
	LEFT JOIN users u ON u.id = e.responsible_person_id
`

func scanEquipment(row interface{ Scan(...any) error }) (*model.Equipment, error) {
	e := &model.Equipment{}
	return e, row.Scan(
		&e.ID, &e.InventoryNumber, &e.Name, &e.Description, &e.CategoryID,
		&e.SerialNumber, &e.Model, &e.Manufacturer, &e.PurchaseDate, &e.PurchasePrice,
		&e.WarrantyExpiry, &e.Status, &e.DepartmentID, &e.ResponsiblePersonID,
		&e.Notes, &e.IsArchived, &e.CreatedAt, &e.UpdatedAt,
		&e.CategoryName, &e.DepartmentName, &e.ResponsiblePersonName,
	)
}

func (r *EquipmentRepository) Create(ctx context.Context, e *model.Equipment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO equipment (id, inventory_number, name, description, category_id,
		    serial_number, model, manufacturer, purchase_date, purchase_price,
		    warranty_expiry, status, department_id, responsible_person_id, notes,
		    is_archived, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
	`, e.ID, e.InventoryNumber, e.Name, e.Description, e.CategoryID,
		e.SerialNumber, e.Model, e.Manufacturer, e.PurchaseDate, e.PurchasePrice,
		e.WarrantyExpiry, e.Status, e.DepartmentID, e.ResponsiblePersonID, e.Notes,
		e.IsArchived, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Equipment, error) {
	row := r.pool.QueryRow(ctx, equipmentSelectBase+" WHERE e.id=$1", id)
	e, err := scanEquipment(row)
	if err != nil {
		return nil, fmt.Errorf("get equipment by id: %w", err)
	}
	photos, err := r.GetPhotos(ctx, id)
	if err != nil {
		return nil, err
	}
	e.Photos = photos
	return e, nil
}

func (r *EquipmentRepository) List(ctx context.Context, f model.EquipmentFilter) ([]*model.Equipment, int, error) {
	where := []string{"e.is_archived = FALSE"}
	args := []any{}
	i := 1

	if f.Search != "" {
		where = append(where, fmt.Sprintf(
			"(e.name ILIKE $%d OR e.inventory_number ILIKE $%d OR e.serial_number ILIKE $%d)",
			i, i, i))
		args = append(args, "%"+f.Search+"%")
		i++
	}
	if f.CategoryID != nil {
		where = append(where, fmt.Sprintf("e.category_id=$%d", i))
		args = append(args, *f.CategoryID)
		i++
	}
	if f.DepartmentID != nil {
		where = append(where, fmt.Sprintf("e.department_id=$%d", i))
		args = append(args, *f.DepartmentID)
		i++
	}
	if f.Status != nil {
		where = append(where, fmt.Sprintf("e.status=$%d", i))
		args = append(args, *f.Status)
		i++
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	var total int
	countSQL := "SELECT COUNT(*) FROM equipment e " + whereSQL
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count equipment: %w", err)
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

	listSQL := equipmentSelectBase + whereSQL +
		fmt.Sprintf(" ORDER BY e.created_at DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list equipment: %w", err)
	}
	defer rows.Close()

	var list []*model.Equipment
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, e)
	}
	return list, total, nil
}

func (r *EquipmentRepository) ListAll(ctx context.Context) ([]*model.Equipment, error) {
	rows, err := r.pool.Query(ctx, equipmentSelectBase+" WHERE e.is_archived = FALSE ORDER BY e.inventory_number")
	if err != nil {
		return nil, fmt.Errorf("list all equipment: %w", err)
	}
	defer rows.Close()
	var list []*model.Equipment
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *EquipmentRepository) ListByDepartment(ctx context.Context, departmentID uuid.UUID) ([]*model.Equipment, error) {
	rows, err := r.pool.Query(ctx,
		equipmentSelectBase+" WHERE e.department_id=$1 AND e.is_archived=FALSE ORDER BY e.inventory_number",
		departmentID)
	if err != nil {
		return nil, fmt.Errorf("list equipment by department: %w", err)
	}
	defer rows.Close()
	var list []*model.Equipment
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, e *model.Equipment) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE equipment SET name=$1, description=$2, category_id=$3, serial_number=$4,
		    model=$5, manufacturer=$6, purchase_date=$7, purchase_price=$8,
		    warranty_expiry=$9, status=$10, department_id=$11, responsible_person_id=$12,
		    notes=$13, updated_at=$14
		WHERE id=$15
	`, e.Name, e.Description, e.CategoryID, e.SerialNumber,
		e.Model, e.Manufacturer, e.PurchaseDate, e.PurchasePrice,
		e.WarrantyExpiry, e.Status, e.DepartmentID, e.ResponsiblePersonID,
		e.Notes, time.Now(), e.ID)
	if err != nil {
		return fmt.Errorf("update equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepository) Archive(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE equipment SET is_archived=TRUE, updated_at=$1 WHERE id=$2",
		time.Now(), id)
	if err != nil {
		return fmt.Errorf("archive equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepository) UpdateDepartment(ctx context.Context, id, departmentID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE equipment SET department_id=$1, updated_at=$2 WHERE id=$3",
		departmentID, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update equipment department: %w", err)
	}
	return nil
}

// Photos

func (r *EquipmentRepository) AddPhoto(ctx context.Context, p *model.EquipmentPhoto) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO equipment_photos (id, equipment_id, file_path, uploaded_at) VALUES ($1,$2,$3,$4)",
		p.ID, p.EquipmentID, p.FilePath, p.UploadedAt)
	if err != nil {
		return fmt.Errorf("add photo: %w", err)
	}
	return nil
}

func (r *EquipmentRepository) GetPhotos(ctx context.Context, equipmentID uuid.UUID) ([]model.EquipmentPhoto, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, equipment_id, file_path, uploaded_at FROM equipment_photos WHERE equipment_id=$1 ORDER BY uploaded_at",
		equipmentID)
	if err != nil {
		return nil, fmt.Errorf("get photos: %w", err)
	}
	defer rows.Close()
	var photos []model.EquipmentPhoto
	for rows.Next() {
		p := model.EquipmentPhoto{}
		if err := rows.Scan(&p.ID, &p.EquipmentID, &p.FilePath, &p.UploadedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *EquipmentRepository) GetPhotoByID(ctx context.Context, photoID uuid.UUID) (*model.EquipmentPhoto, error) {
	p := &model.EquipmentPhoto{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, equipment_id, file_path, uploaded_at FROM equipment_photos WHERE id=$1",
		photoID).Scan(&p.ID, &p.EquipmentID, &p.FilePath, &p.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("get photo by id: %w", err)
	}
	return p, nil
}

func (r *EquipmentRepository) CountPhotos(ctx context.Context, equipmentID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM equipment_photos WHERE equipment_id=$1", equipmentID).Scan(&count)
	return count, err
}

func (r *EquipmentRepository) DeletePhoto(ctx context.Context, photoID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM equipment_photos WHERE id=$1", photoID)
	if err != nil {
		return fmt.Errorf("delete photo: %w", err)
	}
	return nil
}

// Stats

func (r *EquipmentRepository) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT status, COUNT(*) FROM equipment WHERE is_archived=FALSE GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, nil
}

func (r *EquipmentRepository) CountByCategory(ctx context.Context) ([]map[string]any, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.name, COUNT(e.id)
		FROM equipment e
		LEFT JOIN categories c ON c.id = e.category_id
		WHERE e.is_archived=FALSE
		GROUP BY c.name ORDER BY COUNT(e.id) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"category": name, "count": count})
	}
	return result, nil
}

func (r *EquipmentRepository) GetWarrantyExpiringSoon(ctx context.Context, days int) ([]*model.Equipment, error) {
	rows, err := r.pool.Query(ctx,
		equipmentSelectBase+
			" WHERE e.is_archived=FALSE AND e.warranty_expiry BETWEEN NOW() AND NOW() + ($1 * interval '1 day') ORDER BY e.warranty_expiry",
		days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Equipment
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}
