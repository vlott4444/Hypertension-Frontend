package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// Получить все опубликованные карточки.
func (r *HypertensionRepository) PublishedHypertensionServices() ([]ds.HypertensionService, error) {
	var services []ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where("status = ?", ds.HypertensionPublished).
		Order("id ASC").
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

// Получить опубликованную карточку по ID.
// Удаленные и черновые карточки через URL открыть нельзя.
func (r *HypertensionRepository) HypertensionServiceByID(id int) (ds.HypertensionService, error) {
	var service ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where("id = ? AND status = ?", id, ds.HypertensionPublished).
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

// Получить черновик. Если черновика еще нет, он создается через ORM
// при первом открытии страницы /draft.
func (r *HypertensionRepository) HypertensionDraftService() (ds.HypertensionService, error) {
	var service ds.HypertensionService

	err := r.db.
		Preload("Likes").
		Where("status = ?", ds.HypertensionDraft).
		First(&service).
		Error

	if err == nil {
		return service, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionService{}, err
	}

	service = ds.HypertensionService{
		Title:           "Новая карточка стадии гипертонии",
		Description:     "",
		FullDescription: "",
		SystolicBP:      140,
		DiastolicBP:     90,
		Status:          ds.HypertensionDraft,
		// URL намеренно пустые: в HTML предусмотрены локальные фото/видео по умолчанию.
		ImageURL: "",
		VideoURL: "",
	}

	if err := r.db.Create(&service).Error; err != nil {
		return ds.HypertensionService{}, err
	}

	return service, nil
}

// Опубликовать черновик через ORM.
func (r *HypertensionRepository) PublishHypertensionDraft(
	id int,
	description string,
	systolicBP int,
	diastolicBP int,
) error {
	updates := map[string]interface{}{
		"title":        ds.HypertensionStageTitleBySBP(systolicBP),
		"description":  strings.TrimSpace(description),
		"systolic_bp":  systolicBP,
		"diastolic_bp": diastolicBP,
		"status":       ds.HypertensionPublished,
	}

	result := r.db.
		Model(&ds.HypertensionService{}).
		Where("id = ? AND status = ?", id, ds.HypertensionDraft).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("черновик с id=%d не найден", id)
	}

	return nil
}

// Логическое удаление опубликованной карточки ЧИСТЫМ SQL UPDATE, без ORM.
func (r *HypertensionRepository) DeleteHypertensionServiceSQL(
	ctx context.Context,
	id int,
) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	result, err := sqlDB.ExecContext(
		ctx,
		`UPDATE hypertension_services
		 SET status = $1
		 WHERE id = $2 AND status = $3`,
		ds.HypertensionDeleted,
		id,
		ds.HypertensionPublished,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("опубликованная карточка с id=%d не найдена", id)
	}

	return nil
}

// Следующая опубликованная карточка.
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

// Фильтр по систолическому давлению.
func (r *HypertensionRepository) FilterPublishedHypertensionBySBP(
	sbp int,
) ([]ds.HypertensionService, error) {
	stage := ds.HypertensionStageNumberBySBP(sbp)

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
		Order("id ASC").
		Find(&services).
		Error

	if err != nil {
		return nil, err
	}

	return services, nil
}

// Диапазон давления для SQL.
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
