package service

import (
	"context"

	"github.com/siers22/praktika/back/internal/repository"
)

type ReportService struct {
	equipmentRepo *repository.EquipmentRepository
	movementRepo  *repository.MovementRepository
	inventoryRepo *repository.InventoryRepository
}

func NewReportService(equipmentRepo *repository.EquipmentRepository, movementRepo *repository.MovementRepository, inventoryRepo *repository.InventoryRepository) *ReportService {
	return &ReportService{
		equipmentRepo: equipmentRepo,
		movementRepo:  movementRepo,
		inventoryRepo: inventoryRepo,
	}
}

type SummaryReport struct {
	TotalEquipment   int              `json:"total_equipment"`
	ByStatus         map[string]int   `json:"by_status"`
	ByCategory       []map[string]any `json:"by_category"`
}

type DashboardData struct {
	TotalEquipment       int              `json:"total_equipment"`
	ByStatus             map[string]int   `json:"by_status"`
	ByCategory           []map[string]any `json:"by_category"`
	RecentMovements      any              `json:"recent_movements"`
	WarrantyExpiringSoon any              `json:"warranty_expiring_soon"`
}

func (s *ReportService) Summary(ctx context.Context) (*SummaryReport, error) {
	byStatus, err := s.equipmentRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	byCategory, err := s.equipmentRepo.CountByCategory(ctx)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, v := range byStatus {
		total += v
	}

	return &SummaryReport{
		TotalEquipment: total,
		ByStatus:       byStatus,
		ByCategory:     byCategory,
	}, nil
}

func (s *ReportService) ByDepartment(ctx context.Context) (any, error) {
	list, err := s.equipmentRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	type depStat struct {
		Name          string  `json:"name"`
		Count         int     `json:"count"`
		TotalValue    float64 `json:"total_value"`
	}

	byDep := map[string]*depStat{}
	for _, e := range list {
		name := "Не указано"
		if e.DepartmentName != nil {
			name = *e.DepartmentName
		}
		if _, ok := byDep[name]; !ok {
			byDep[name] = &depStat{Name: name}
		}
		byDep[name].Count++
		if e.PurchasePrice != nil {
			byDep[name].TotalValue += *e.PurchasePrice
		}
	}

	result := make([]*depStat, 0, len(byDep))
	for _, v := range byDep {
		result = append(result, v)
	}
	return result, nil
}

func (s *ReportService) Dashboard(ctx context.Context) (*DashboardData, error) {
	byStatus, err := s.equipmentRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	byCategory, err := s.equipmentRepo.CountByCategory(ctx)
	if err != nil {
		return nil, err
	}

	recentMovements, err := s.movementRepo.GetRecent(ctx, 10)
	if err != nil {
		return nil, err
	}

	warrantySoon, err := s.equipmentRepo.GetWarrantyExpiringSoon(ctx, 30)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, v := range byStatus {
		total += v
	}

	return &DashboardData{
		TotalEquipment:       total,
		ByStatus:             byStatus,
		ByCategory:           byCategory,
		RecentMovements:      recentMovements,
		WarrantyExpiringSoon: warrantySoon,
	}, nil
}
