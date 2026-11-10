package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"hypertensionstages/internal/app/ds"

	"crypto/subtle"
	"golang.org/x/crypto/bcrypt"
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
func (r *HypertensionRepository) PublishedHypertensionStages() ([]ds.HypertensionStage, error) {
	var services []ds.HypertensionStage

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
func (r *HypertensionRepository) HypertensionStageByID(id int) (ds.HypertensionStage, error) {
	var service ds.HypertensionStage

	err := r.db.
		Preload("Likes").
		Where("id = ? AND status = ?", id, ds.HypertensionPublished).
		First(&service).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionStage{}, fmt.Errorf(
			"опубликованная карточка стадии гипертонии с id=%d не найдена",
			id,
		)
	}

	if err != nil {
		return ds.HypertensionStage{}, err
	}

	return service, nil
}

// Получить черновик. Если черновика еще нет, он создается через ORM
// при первом открытии страницы /draft.
func (r *HypertensionRepository) HypertensionDraftService() (ds.HypertensionStage, error) {
	var service ds.HypertensionStage

	err := r.db.
		Preload("Likes").
		Where("status = ?", ds.HypertensionDraft).
		First(&service).
		Error

	if err == nil {
		return service, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionStage{}, err
	}

	service = ds.HypertensionStage{
		Title:       "Новая карточка стадии гипертонии",
		Description: "",
		SystolicBP:  140,
		DiastolicBP: 90,
		Status:      ds.HypertensionDraft,
		// Медиа черновика по умолчанию находятся на сервере в resources/media.
		ImageName: "draft_stage1.png",
		VideoName: "draft_stage1.mp4",
	}

	if err := r.db.Create(&service).Error; err != nil {
		return ds.HypertensionStage{}, err
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
		Model(&ds.HypertensionStage{}).
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

func (r *HypertensionRepository) DeleteHypertensionStageSQL(
	ctx context.Context,
	id int,
) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	result, err := sqlDB.ExecContext(
		ctx,
		`UPDATE hypertension_stages
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
func (r *HypertensionRepository) NextPublishedHypertensionStage(
	currentID int,
) (ds.HypertensionStage, error) {
	services, err := r.PublishedHypertensionStages()
	if err != nil {
		return ds.HypertensionStage{}, err
	}

	for index, service := range services {
		if int(service.ID) == currentID {
			next := services[(index+1)%len(services)]
			return next, nil
		}
	}

	return ds.HypertensionStage{}, fmt.Errorf(
		"невозможно открыть следующую карточку после id=%d",
		currentID,
	)
}

// Фильтр по систолическому давлению.
// Оставлен для обратной совместимости, если где-то ещё вызывается.
func (r *HypertensionRepository) FilterPublishedHypertensionBySBP(
	sbp int,
) ([]ds.HypertensionStage, error) {
	stage := ds.HypertensionStageNumberBySBP(sbp)

	min, max := hypertensionRange(stage)

	var services []ds.HypertensionStage

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

// SBPBounds возвращает минимальное и максимальное систолическое давление
// среди опубликованных карточек. Нужно для границ шкалы двух слайдеров.
func (r *HypertensionRepository) SBPBounds() (int, int, error) {
	var bounds struct {
		Min int
		Max int
	}

	err := r.db.
		Model(&ds.HypertensionStage{}).
		Where("status = ?", ds.HypertensionPublished).
		Select("COALESCE(MIN(systolic_bp), 0) AS min, COALESCE(MAX(systolic_bp), 0) AS max").
		Scan(&bounds).
		Error

	if err != nil {
		return 0, 0, err
	}

	return bounds.Min, bounds.Max, nil
}

// FilterPublishedHypertensionBySBPRange возвращает опубликованные карточки,
// у которых SystolicBP входит в диапазон [min, max].
func (r *HypertensionRepository) FilterPublishedHypertensionBySBPRange(
	min, max int,
) ([]ds.HypertensionStage, error) {
	var services []ds.HypertensionStage

	err := r.db.
		Preload("Likes").
		Where(
			"status = ? AND systolic_bp BETWEEN ? AND ?",
			ds.HypertensionPublished,
			min,
			max,
		).
		Order("systolic_bp ASC").
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

// ============================================================
// API-МЕТОДЫ
// ============================================================

func (r *HypertensionRepository) APIGetStages(minBP, maxBP int, title string) ([]ds.HypertensionStage, error) {
	var stages []ds.HypertensionStage
	q := r.db.Preload("Likes").Where("status = ? AND systolic_bp BETWEEN ? AND ?", ds.HypertensionPublished, minBP, maxBP)
	if title != "" {
		q = q.Where("title ILIKE ?", "%"+title+"%")
	}
	err := q.Order("id ASC").Find(&stages).Error
	return stages, err
}

// APIGetFeed — лента опубликованных.
func (r *HypertensionRepository) APIGetFeed() ([]ds.HypertensionStage, error) {
	var services []ds.HypertensionStage
	err := r.db.
		Preload("Likes").
		Where("status = ?", ds.HypertensionPublished).
		Order("id ASC").
		Find(&services).Error
	return services, err
}

// APIPublishedByID — опубликованная по ID.
func (r *HypertensionRepository) APIPublishedByID(
	id int,
) (ds.HypertensionStage, error) {
	var s ds.HypertensionStage
	err := r.db.
		Preload("Likes").
		Where("id = ? AND status = ?", id, ds.HypertensionPublished).
		First(&s).Error
	return s, err
}

// APINextPublished — следующая опубликованная после ID (по кругу).
func (r *HypertensionRepository) APINextPublished(
	id int,
) (ds.HypertensionStage, error) {
	services, err := r.APIGetFeed()
	if err != nil {
		return ds.HypertensionStage{}, err
	}
	if len(services) == 0 {
		return ds.HypertensionStage{}, fmt.Errorf("no published stages")
	}
	for i, s := range services {
		if int(s.ID) == id {
			next := services[(i+1)%len(services)]
			return next, nil
		}
	}
	return ds.HypertensionStage{}, fmt.Errorf("id=%d not found", id)
}

// APIGetDraft — черновик пользователя.
func (r *HypertensionRepository) APIGetDraft(
	userID uint,
) (ds.HypertensionStage, error) {
	var s ds.HypertensionStage
	err := r.db.
		Where("user_id = ? AND status = ?", userID, ds.HypertensionDraft).
		First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.HypertensionStage{}, nil
	}
	return s, err
}

// APICreateStage — создание услуги.
//
// ВАЖНО: moderator_id обязателен по методичке, и на него есть FK
// на таблицу users. На этапе лабы модератором ставим самого создателя.
func (r *HypertensionRepository) APICreateStage(
	userID uint,
	title, description string,
	sbp, dbp int,
	imageName, videoName string,
) (ds.HypertensionStage, error) {

	s := ds.HypertensionStage{
		Title:       title,
		Description: description,
		SystolicBP:  sbp,
		DiastolicBP: dbp,
		Status:      ds.HypertensionDraft,
		UserID:      userID,
		ModeratorID: userID,
		ImageName:   imageName,
		VideoName:   videoName,
	}

	if err := r.db.Create(&s).Error; err != nil {
		return ds.HypertensionStage{}, err
	}
	return s, nil
}

// APIPublishStage — публикация черновика пользователя.
func (r *HypertensionRepository) APIPublishStage(
	id int, userID uint,
) (ds.HypertensionStage, error) {

	now := time.Now()

	result := r.db.
		Model(&ds.HypertensionStage{}).
		Where("id = ? AND user_id = ? AND status = ?",
			id, userID, ds.HypertensionDraft).
		Updates(map[string]interface{}{
			"status":    ds.HypertensionPublished,
			"formed_at": now,
		})

	if result.Error != nil {
		return ds.HypertensionStage{}, result.Error
	}
	if result.RowsAffected == 0 {
		return ds.HypertensionStage{}, fmt.Errorf("черновик не найден или не ваш")
	}

	var s ds.HypertensionStage
	r.db.First(&s, id)
	return s, nil
}

// APISoftDeleteStage — soft delete только своих.
func (r *HypertensionRepository) APISoftDeleteStage(
	id int, userID uint,
) error {

	result := r.db.
		Model(&ds.HypertensionStage{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", ds.HypertensionDeleted)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("услуга не найдена или не ваша")
	}
	return nil
}

// APIToggleLike — поставить/снять лайк.
func (r *HypertensionRepository) APIToggleLike(
	stageID int, userID uint, value int,
) (ds.HypertensionLike, error) {

	var present int64
	if err := r.db.Model(&ds.HypertensionStage{}).Where("id = ? AND status = ?", stageID, ds.HypertensionPublished).Count(&present).Error; err != nil {
		return ds.HypertensionLike{}, err
	}
	if present == 0 {
		return ds.HypertensionLike{}, fmt.Errorf("stage not found or unpublished")
	}
	var like ds.HypertensionLike
	err := r.db.
		Where("user_id = ? AND hypertension_stage_id = ?", userID, stageID).
		First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		like = ds.HypertensionLike{
			UserID:              userID,
			HypertensionStageID: uint(stageID),
			Value:               value,
		}
		return like, r.db.Create(&like).Error
	}
	if err != nil {
		return ds.HypertensionLike{}, err
	}

	like.Value = value
	return like, r.db.Save(&like).Error
}

// User passwords are stored as bcrypt hashes. Legacy plain-text seed records are
// upgraded the first time they are used for a successful login.
func (r *HypertensionRepository) APIRegisterUser(username, password string) (ds.User, error) {
	var existing int64
	if err := r.db.Model(&ds.User{}).Where("username = ?", username).Count(&existing).Error; err != nil {
		return ds.User{}, err
	}
	if existing != 0 {
		return ds.User{}, fmt.Errorf("username taken")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ds.User{}, err
	}
	user := ds.User{Username: username, Password: string(hashed)}
	if err := r.db.Create(&user).Error; err != nil {
		return ds.User{}, err
	}
	return user, nil
}
func (r *HypertensionRepository) APIAuthenticateUser(username, password string) (ds.User, error) {
	var user ds.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return ds.User{}, fmt.Errorf("invalid credentials")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		return user, nil
	}
	// Existing lab4 data may contain plain-text test passwords; upgrade on login.
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) == 1 {
		hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			return ds.User{}, e
		}
		if e = r.db.Model(&user).Update("password", string(hash)).Error; e != nil {
			return ds.User{}, e
		}
		return user, nil
	}
	return ds.User{}, fmt.Errorf("invalid credentials")
}
