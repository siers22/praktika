package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/repository"
)

type EquipmentService struct {
	repo      *repository.EquipmentRepository
	auditRepo *repository.AuditRepository
	uploadDir string
}

func NewEquipmentService(repo *repository.EquipmentRepository, auditRepo *repository.AuditRepository, uploadDir string) *EquipmentService {
	return &EquipmentService{repo: repo, auditRepo: auditRepo, uploadDir: uploadDir}
}

func (s *EquipmentService) Create(ctx context.Context, actorID uuid.UUID, req *model.CreateEquipmentRequest) (*model.Equipment, error) {
	now := time.Now()
	e := &model.Equipment{
		ID:                  uuid.New(),
		InventoryNumber:     req.InventoryNumber,
		Name:                req.Name,
		Description:         req.Description,
		CategoryID:          req.CategoryID,
		SerialNumber:        req.SerialNumber,
		Model:               req.Model,
		Manufacturer:        req.Manufacturer,
		PurchaseDate:        req.PurchaseDate,
		PurchasePrice:       req.PurchasePrice,
		WarrantyExpiry:      req.WarrantyExpiry,
		Status:              req.Status,
		DepartmentID:        req.DepartmentID,
		ResponsiblePersonID: req.ResponsiblePersonID,
		Notes:               req.Notes,
		IsArchived:          false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, actorID, "create", "equipment", &e.ID, map[string]any{"inventory_number": e.InventoryNumber})
	return s.repo.GetByID(ctx, e.ID)
}

func (s *EquipmentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Equipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EquipmentService) List(ctx context.Context, f model.EquipmentFilter) ([]*model.Equipment, int, error) {
	if f.PerPage <= 0 || f.PerPage > 50 {
		f.PerPage = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	return s.repo.List(ctx, f)
}

func (s *EquipmentService) Update(ctx context.Context, actorID, id uuid.UUID, req *model.UpdateEquipmentRequest) (*model.Equipment, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.IsArchived {
		return nil, errors.New("cannot update archived equipment")
	}

	e.Name = req.Name
	e.Description = req.Description
	e.CategoryID = req.CategoryID
	e.SerialNumber = req.SerialNumber
	e.Model = req.Model
	e.Manufacturer = req.Manufacturer
	e.PurchaseDate = req.PurchaseDate
	e.PurchasePrice = req.PurchasePrice
	e.WarrantyExpiry = req.WarrantyExpiry
	e.Status = req.Status
	e.DepartmentID = req.DepartmentID
	e.ResponsiblePersonID = req.ResponsiblePersonID
	e.Notes = req.Notes

	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, actorID, "update", "equipment", &id, req)
	return s.repo.GetByID(ctx, id)
}

func (s *EquipmentService) Archive(ctx context.Context, actorID, id uuid.UUID) error {
	if err := s.repo.Archive(ctx, id); err != nil {
		return err
	}
	s.auditRepo.Log(ctx, actorID, "archive", "equipment", &id, nil)
	return nil
}

func (s *EquipmentService) UploadPhoto(ctx context.Context, actorID, equipmentID uuid.UUID, fh *multipart.FileHeader) (*model.EquipmentPhoto, error) {
	count, err := s.repo.CountPhotos(ctx, equipmentID)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, errors.New("maximum 5 photos allowed per equipment")
	}

	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	dir := filepath.Join(s.uploadDir, equipmentID.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}

	photoID := uuid.New()
	ext := filepath.Ext(fh.Filename)
	filename := photoID.String() + ext
	dstPath := filepath.Join(dir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	relativePath := filepath.Join(equipmentID.String(), filename)
	photo := &model.EquipmentPhoto{
		ID:          photoID,
		EquipmentID: equipmentID,
		FilePath:    relativePath,
		UploadedAt:  time.Now(),
	}

	if err := s.repo.AddPhoto(ctx, photo); err != nil {
		_ = os.Remove(dstPath)
		return nil, err
	}

	slog.Info("photo uploaded", "equipment_id", equipmentID, "photo_id", photoID)
	return photo, nil
}

func (s *EquipmentService) DeletePhoto(ctx context.Context, actorID, equipmentID, photoID uuid.UUID) error {
	photo, err := s.repo.GetPhotoByID(ctx, photoID)
	if err != nil {
		return err
	}
	if photo.EquipmentID != equipmentID {
		return errors.New("photo does not belong to this equipment")
	}

	_ = os.Remove(filepath.Join(s.uploadDir, photo.FilePath))
	return s.repo.DeletePhoto(ctx, photoID)
}

func (s *EquipmentService) ExportCSV(ctx context.Context, w io.Writer) error {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return err
	}

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{
		"Инвентарный номер", "Наименование", "Категория", "Серийный номер",
		"Производитель", "Модель", "Статус", "Подразделение",
		"Дата приобретения", "Стоимость", "Окончание гарантии",
	})

	for _, e := range list {
		purchaseDate := ""
		if e.PurchaseDate != nil {
			purchaseDate = e.PurchaseDate.Format("2006-01-02")
		}
		price := ""
		if e.PurchasePrice != nil {
			price = strconv.FormatFloat(*e.PurchasePrice, 'f', 2, 64)
		}
		warranty := ""
		if e.WarrantyExpiry != nil {
			warranty = e.WarrantyExpiry.Format("2006-01-02")
		}
		cat := ""
		if e.CategoryName != nil {
			cat = *e.CategoryName
		}
		dep := ""
		if e.DepartmentName != nil {
			dep = *e.DepartmentName
		}
		serial := ""
		if e.SerialNumber != nil {
			serial = *e.SerialNumber
		}
		manufacturer := ""
		if e.Manufacturer != nil {
			manufacturer = *e.Manufacturer
		}
		model := ""
		if e.Model != nil {
			model = *e.Model
		}
		_ = cw.Write([]string{
			e.InventoryNumber, e.Name, cat, serial,
			manufacturer, model, string(e.Status), dep,
			purchaseDate, price, warranty,
		})
	}
	cw.Flush()
	return cw.Error()
}
