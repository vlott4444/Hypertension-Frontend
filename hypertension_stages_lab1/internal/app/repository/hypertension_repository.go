package repository

import (
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

// Получить все опубликованные карточки гипертонии
func (r *HypertensionRepository) PublishedHypertensionServices() ([]ds.HypertensionService, error) {
	var services []ds.HypertensionService

	err := r.db.
		Where("status = ?", ds.HypertensionPublished).
		Find(&services).
		Error

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("опубликованные карточки стадий гипертонии не найдены")
	}

	return services, nil
}

// Получить карточку по ID
func (r *HypertensionRepository) HypertensionServiceByID(id int) (ds.HypertensionService, error) {
	var service ds.HypertensionService

	err := r.db.
		Where("id = ? AND status = ?", id, ds.HypertensionPublished).
		First(&service).
		Error

	if err != nil {
		return ds.HypertensionService{}, fmt.Errorf(
			"опубликованная карточка стадии гипертонии с id=%d не найдена",
			id,
		)
	}

	return service, nil
}

// Получить черновик
func (r *HypertensionRepository) HypertensionDraftService() (ds.HypertensionService, error) {
	var service ds.HypertensionService

	err := r.db.
		Where("status = ?", ds.HypertensionDraft).
		First(&service).
		Error

	if err != nil {
		return ds.HypertensionService{}, fmt.Errorf(
			"черновик стадии гипертонии не найден",
		)
	}

	return service, nil
}

// Следующая опубликованная карточка
func (r *HypertensionRepository) NextPublishedHypertensionService(currentID int) (ds.HypertensionService, error) {

	services, err := r.PublishedHypertensionServices()

	if err != nil {
		return ds.HypertensionService{}, err
	}

	for index, service := range services {

		if service.ID == currentID {

			nextIndex := (index + 1) % len(services)

			return services[nextIndex], nil
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

	services, err := r.PublishedHypertensionServices()

	if err != nil {
		return nil, err
	}

	result := make([]ds.HypertensionService, 0)

	for _, service := range services {

		if ds.HypertensionStageNumberBySBP(service.SystolicBP) == stage {
			result = append(result, service)
		}
	}

	return result, nil
}
