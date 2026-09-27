package ds

type User struct {
	ID uint `gorm:"primaryKey"`

	Username string

	Password string

	Likes []HypertensionLike
}
