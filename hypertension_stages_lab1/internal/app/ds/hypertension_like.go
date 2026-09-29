package ds

type HypertensionLike struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	HypertensionServiceID uint                `json:"hypertension_service_id"`
	HypertensionService   HypertensionService `gorm:"foreignKey:HypertensionServiceID" json:"-"`

	// 0 — лайк снят, 1 — лайк поставлен
	Value int `json:"value"`
}
