package repository

import (
	"errors"
	"fmt"

	"hypertensionstages/internal/app/ds"

	"gorm.io/gorm"
)

type HypertensionRepository struct {
	db *gorm.DB
}

func NewHypertensionRepository(db *gorm.DB) *HypertensionRepository {
	return &HypertensionRepository{
		db: db,
	}
}

// Получить все опубликованные карточки
func (r *HypertensionRepository) PublishedHypertensionServices() ([]ds.HypertensionService, error) {

	var services []ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where("status = ?", ds.HypertensionPublished).
		Find(&services).
		Error

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf(
			"опубликованные карточки стадий гипертонии не найдены",
		)
	}

	return services, nil
}

// Получить карточку по ID
func (r *HypertensionRepository) HypertensionServiceByID(id int) (ds.HypertensionService, error) {

	var service ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where(
			"id = ? AND status = ?",
			id,
			ds.HypertensionPublished,
		).
		First(&service).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionService{}, fmt.Errorf(
			"опубликованная карточка стадии гипертонии с id=%d не найдена",
			id,
		)
	}

	if err != nil {
		return ds.HypertensionService{}, err
	}

	return service, nil
}

// Получить черновик
func (r *HypertensionRepository) HypertensionDraftService() (ds.HypertensionService, error) {

	var service ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where(
			"status = ?",
			ds.HypertensionDraft,
		).
		First(&service).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionService{}, fmt.Errorf(
			"черновик стадии гипертонии не найден",
		)
	}

	if err != nil {
		return ds.HypertensionService{}, err
	}

	return service, nil
}

// Следующая опубликованная карточка
func (r *HypertensionRepository) NextPublishedHypertensionService(
	currentID int,
) (ds.HypertensionService, error) {

	services, err := r.PublishedHypertensionServices()

	if err != nil {
		return ds.HypertensionService{}, err
	}

	for index, service := range services {

		if int(service.ID) == currentID {

			next := services[(index+1)%len(services)]

			return next, nil
		}
	}

	return ds.HypertensionService{}, fmt.Errorf(
		"невозможно открыть следующую карточку после id=%d",
		currentID,
	)
}

// Фильтр по систолическому давлению
func (r *HypertensionRepository) FilterPublishedHypertensionBySBP(
	sbp int,
) ([]ds.HypertensionService, error) {

	stage := ds.HypertensionStageNumberBySBP(sbp)

	if stage == 99 {
		return []ds.HypertensionService{}, nil
	}

	min, max := hypertensionRange(stage)

	var services []ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where(
			"status = ? AND systolic_bp BETWEEN ? AND ?",
			ds.HypertensionPublished,
			min,
			max,
		).
		Find(&services).
		Error

	if err != nil {
		return nil, err
	}

	return services, nil
}

// Диапазон давления для SQL
func hypertensionRange(stage int) (int, int) {

	switch stage {

	case -2:
		return 0, 119

	case -1:
		return 120, 129

	case 0:
		return 130, 139

	case 1:
		return 140, 159

	case 2:
		return 160, 179

	default:
		return 180, 300
	}
}
