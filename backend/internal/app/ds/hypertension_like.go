package ds

type HypertensionLike struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"uniqueIndex:idx_stage_user" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	HypertensionStageID uint              `gorm:"uniqueIndex:idx_stage_user" json:"hypertension_stage_id"`
	HypertensionStage   HypertensionStage `gorm:"foreignKey:HypertensionStageID" json:"-"`

	// 0 — лайк снят, 1 — лайк поставлен
	Value int `json:"value"`
}
