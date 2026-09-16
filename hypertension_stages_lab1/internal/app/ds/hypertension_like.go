package ds

type HypertensionLike struct {
	ID                    uint `gorm:"primaryKey"`
	UserID                uint
	User                  User `gorm:"foreignKey:UserID"`
	HypertensionServiceID uint
	HypertensionService   HypertensionService `gorm:"foreignKey:HypertensionServiceID"`
}
