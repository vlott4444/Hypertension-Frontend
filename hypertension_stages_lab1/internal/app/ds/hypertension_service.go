package ds

type HypertensionStatus string

const (
	HypertensionDraft     HypertensionStatus = "черновик"
	HypertensionPublished HypertensionStatus = "опубликован"
	HypertensionDeleted   HypertensionStatus = "удален"
)

type HypertensionService struct {
	ID              uint `gorm:"primaryKey"`
	Title           string
	Description     string
	FullDescription string
	SystolicBP      int
	DiastolicBP     int
	Status          HypertensionStatus
	ImageURL        string
	VideoURL        string

	Likes []HypertensionLike `gorm:"foreignKey:HypertensionServiceID"`
}

// --------------------
// Логика определения стадии
// --------------------

func HypertensionStageNumberBySBP(sbp int) int {

	switch {

	case sbp < 120:
		return -2

	case sbp <= 129:
		return -1

	case sbp <= 139:
		return 0

	case sbp <= 159:
		return 1

	case sbp <= 179:
		return 2

	default:
		return 3
	}
}

func HypertensionStageRomanTitleBySBP(sbp int) string {

	switch HypertensionStageNumberBySBP(sbp) {

	case 1:
		return "I стадия"

	case 2:
		return "II стадия"

	case 3:
		return "III стадия"

	default:
		return "Вне диапазона"
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
