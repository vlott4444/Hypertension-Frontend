package ds

// CurrentUserID — константа, фиксирующий пользователя-создателя.
// Используется во всех методах вместо реальной авторизации.
const CurrentUserID uint = 1

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Username string `json:"username"`

	// Пароль не отдаём в JSON
	Password string `json:"-"`

	Likes []HypertensionLike `json:"-"`
}

// GetCurrentUser — функция-singleton.
// Возвращает фиксированного пользователя-создателя.
// Используется во всех хендлерах вместо авторизации.
func GetCurrentUser() *User {
	return &User{
		ID:       CurrentUserID,
		Username: "test_user",
	}
}
