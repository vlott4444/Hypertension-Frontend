package repository

import "fmt"

type HypertensionStatus string

const (
	HypertensionDraft     HypertensionStatus = "черновик"
	HypertensionPublished HypertensionStatus = "опубликован"
	HypertensionDeleted   HypertensionStatus = "удален"
)

type HypertensionLike struct {
	UserID int
}

type HypertensionService struct {
	ID          int
	Title       string
	Description string
	SystolicBP  int
	DiastolicBP int
	Status      HypertensionStatus
	ImageURL    string
	VideoURL    string
	Likes       []HypertensionLike
}

type HypertensionRepository struct {
	hypertensionServices []HypertensionService
}

func NewHypertensionRepository() (*HypertensionRepository, error) {
	const minioBaseURL = "http://localhost:9000/hypertension-media"

	hypertensionServices := []HypertensionService{
		{
			ID:          6,
			Title:       "Оптимальное давление",
			Description: "Систолическое давление ниже 120 мм рт. ст., диастолическое ниже 80 мм рт. ст.",
			SystolicBP:  110,
			DiastolicBP: 70,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage_optimal.png",
			VideoURL:    minioBaseURL + "/optimal.mp4",
			Likes:       []HypertensionLike{},
		},
		{
			ID:          7,
			Title:       "Нормальное давление",
			Description: "Систолическое давление 120–129 мм рт. ст. и/или диастолическое 80–84 мм рт. ст.",
			SystolicBP:  125,
			DiastolicBP: 82,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage_normal.png",
			VideoURL:    minioBaseURL + "/normal.mp4",
			Likes:       []HypertensionLike{},
		},
		{
			ID:          8,
			Title:       "Высокое нормальное давление",
			Description: "Систолическое давление 130–139 мм рт. ст. и/или диастолическое 85–89 мм рт. ст.",
			SystolicBP:  135,
			DiastolicBP: 87,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage_high_normal.png",
			VideoURL:    minioBaseURL + "/high_normal.mp4",
			Likes:       []HypertensionLike{},
		},
		{
			ID:          1,
			Title:       "1 стадия",
			Description: "Систолическое артериальное давление находится в диапазоне 140–159 мм рт. ст.",
			SystolicBP:  150,
			DiastolicBP: 95,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage1.png",
			VideoURL:    minioBaseURL + "/stage1.mp4",
			Likes:       []HypertensionLike{{UserID: 101}, {UserID: 205}},
		},
		{
			ID:          2,
			Title:       "2 стадия",
			Description: "Систолическое артериальное давление находится в диапазоне 160–179 мм рт. ст.",
			SystolicBP:  170,
			DiastolicBP: 105,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage2.png",
			VideoURL:    minioBaseURL + "/stage2.mp4",
			Likes:       []HypertensionLike{{UserID: 101}, {UserID: 310}, {UserID: 411}},
		},
		{
			ID:          3,
			Title:       "3 стадия",
			Description: "Систолическое артериальное давление составляет 180 мм рт. ст. или выше.",
			SystolicBP:  190,
			DiastolicBP: 115,
			Status:      HypertensionPublished,
			ImageURL:    minioBaseURL + "/stage3.png",
			VideoURL:    minioBaseURL + "/stage3.mp4",
			Likes:       []HypertensionLike{{UserID: 205}},
		},
		{
			ID:          4,
			Title:       "1 стадия — черновик",
			Description: "Черновик карточки для страницы добавления.",
			SystolicBP:  155,
			DiastolicBP: 98,
			Status:      HypertensionDraft,
			ImageURL:    minioBaseURL + "/draft_stage1.png",
			VideoURL:    minioBaseURL + "/draft_stage1.mp4",
			Likes:       []HypertensionLike{},
		},
		{
			ID:          5,
			Title:       "2 стадия — удаленная карточка",
			Description: "Эта услуга имеет статус удален и намеренно не показывается в интерфейсе.",
			SystolicBP:  165,
			DiastolicBP: 102,
			Status:      HypertensionDeleted,
			ImageURL:    minioBaseURL + "/deleted_stage2.png",
			VideoURL:    minioBaseURL + "/deleted_stage2.mp4",
			Likes:       []HypertensionLike{{UserID: 999}},
		},
	}

	if len(hypertensionServices) == 0 {
		return nil, fmt.Errorf("коллекция стадий гипертонической болезни пуста")
	}

	return &HypertensionRepository{hypertensionServices: hypertensionServices}, nil
}

func (r *HypertensionRepository) PublishedHypertensionServices() ([]HypertensionService, error) {
	result := make([]HypertensionService, 0)
	for _, hypertensionService := range r.hypertensionServices {
		if hypertensionService.Status == HypertensionPublished {
			result = append(result, hypertensionService)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("опубликованные карточки стадий гипертонии не найдены")
	}
	return result, nil
}

func (r *HypertensionRepository) HypertensionServiceByID(id int) (HypertensionService, error) {
	for _, hypertensionService := range r.hypertensionServices {
		if hypertensionService.ID == id && hypertensionService.Status == HypertensionPublished {
			return hypertensionService, nil
		}
	}
	return HypertensionService{}, fmt.Errorf("опубликованная карточка стадии гипертонии с id=%d не найдена", id)
}

func (r *HypertensionRepository) HypertensionDraftService() (HypertensionService, error) {
	for _, hypertensionService := range r.hypertensionServices {
		if hypertensionService.Status == HypertensionDraft {
			return hypertensionService, nil
		}
	}
	return HypertensionService{}, fmt.Errorf("черновик стадии гипертонии не найден")
}

func (r *HypertensionRepository) NextPublishedHypertensionService(currentID int) (HypertensionService, error) {
	published, err := r.PublishedHypertensionServices()
	if err != nil {
		return HypertensionService{}, err
	}

	for index, hypertensionService := range published {
		if hypertensionService.ID == currentID {
			return published[(index+1)%len(published)], nil
		}
	}
	return HypertensionService{}, fmt.Errorf("невозможно открыть следующую карточку после id=%d", currentID)
}

func (r *HypertensionRepository) FilterPublishedHypertensionBySBP(sbp int) ([]HypertensionService, error) {
	requestedStage := HypertensionStageNumberBySBP(sbp)
	if requestedStage == 99 {
		return []HypertensionService{}, nil
	}

	published, err := r.PublishedHypertensionServices()
	if err != nil {
		return nil, err
	}

	result := make([]HypertensionService, 0)
	for _, hypertensionService := range published {
		if HypertensionStageNumberBySBP(hypertensionService.SystolicBP) == requestedStage {
			result = append(result, hypertensionService)
		}
	}
	return result, nil
}

func HypertensionStageNumberBySBP(sbp int) int {
	switch {
	case sbp < 120:
		return -2 // оптимальное
	case sbp <= 129:
		return -1 // нормальное
	case sbp <= 139:
		return 0 // высокое нормальное
	case sbp <= 159:
		return 1
	case sbp <= 179:
		return 2
	default:
		return 99
	}
}

func HypertensionStageRomanTitleBySBP(sbp int) string {
	switch HypertensionStageNumberBySBP(sbp) {
	case 1:
		return "I стадия"
	case 2:
		return "II стадия"
	case 0:
		return "Вне диапазона"
	case -1:
		return "Вне диапазона"
	case -2:
		return "Вне диапазона"
	default:
		return "III стадия"
	}
}

func HypertensionStageTitleBySBP(sbp int) string {
	switch HypertensionStageNumberBySBP(sbp) {
	case -2:
		return "Оптимальное давление"
	case -1:
		return "Нормальное давление"
	case 0:
		return "Высокое нормальное давление"
	case 1:
		return "1 стадия"
	case 2:
		return "2 стадия"
	case 3:
		return "3 стадия"
	default:
		return "3 стадия"
	}
}

func HypertensionStageRangeBySBP(sbp int) string {
	switch {
	case sbp < 120:
		return "<120 мм рт. ст."
	case sbp <= 129:
		return "120–129 мм рт. ст."
	case sbp <= 139:
		return "130–139 мм рт. ст."
	case sbp <= 159:
		return "140–159 мм рт. ст."
	case sbp <= 179:
		return "160–179 мм рт. ст."
	default:
		return "≥180 мм рт. ст."
	}
}
