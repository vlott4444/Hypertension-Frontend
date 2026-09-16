package main

import (
	"hypertensionstages/internal/app/ds"
	"hypertensionstages/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	services := []ds.HypertensionService{

		{
			Title: "Оптимальное давление",

			Description: "Систолическое давление ниже 120 мм рт. ст.",

			FullDescription: "Оптимальный уровень артериального давления характеризуется значениями ниже 120 мм рт. ст. для систолического давления и ниже 80 мм рт. ст. для диастолического. Обычно такое давление связано с низким риском сердечно-сосудистых осложнений. Для поддержания нормальных показателей рекомендуется сохранять физическую активность, придерживаться сбалансированного питания и регулярно контролировать давление.",

			SystolicBP:  110,
			DiastolicBP: 70,

			ImageURL: "/resources/media/stage_optimal.png",
			VideoURL: "/resources/media/optimal.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "Нормальное давление",

			Description: "Систолическое давление 120–129 мм рт. ст.",

			FullDescription: "Нормальное артериальное давление находится в диапазоне 120–129 мм рт. ст. по систолическому показателю. Такие значения считаются благоприятными, однако важно следить за динамикой давления. Регулярные измерения помогают вовремя заметить изменения и предотвратить развитие стойкой гипертонии.",

			SystolicBP:  125,
			DiastolicBP: 78,

			ImageURL: "/resources/media/stage_normal.png",
			VideoURL: "/resources/media/normal.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "Высокое нормальное давление",

			Description: "Систолическое давление 130–139 мм рт. ст.",

			FullDescription: "Повышенное нормальное давление означает, что показатели уже находятся выше оптимального уровня. Это состояние не всегда требует лечения, но является сигналом для изменения образа жизни. Важно контролировать массу тела, уровень физической активности, питание и регулярно измерять давление.",

			SystolicBP:  135,
			DiastolicBP: 85,

			ImageURL: "/resources/media/stage_high_normal.png",
			VideoURL: "/resources/media/high_normal.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "I стадия гипертонии",

			Description: "Систолическое давление 140–159 мм рт. ст.",

			FullDescription: "Первая стадия артериальной гипертонии характеризуется устойчивым повышением давления. На этом этапе могут отсутствовать выраженные симптомы, поэтому регулярный контроль особенно важен. Врач оценивает общий риск осложнений и при необходимости назначает лечение вместе с рекомендациями по изменению образа жизни.",

			SystolicBP:  150,
			DiastolicBP: 95,

			ImageURL: "/resources/media/stage1.png",
			VideoURL: "/resources/media/stage1.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "II стадия гипертонии",

			Description: "Систолическое давление 160–179 мм рт. ст.",

			FullDescription: "Вторая стадия гипертонии характеризуется более значительным повышением давления. При длительном сохранении таких значений увеличивается нагрузка на сердце, сосуды и другие органы. Обычно требуется регулярное наблюдение специалиста и подбор терапии для достижения целевых показателей давления.",

			SystolicBP:  170,
			DiastolicBP: 105,

			ImageURL: "/resources/media/stage2.png",
			VideoURL: "/resources/media/stage2.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "III стадия гипертонии",

			Description: "Систолическое давление 180 мм рт. ст. и выше",

			FullDescription: "Третья стадия артериальной гипертонии характеризуется выраженным повышением давления. Такие значения требуют особого внимания, поскольку связаны с повышенной нагрузкой на сердечно-сосудистую систему. Необходимо регулярное медицинское наблюдение и соблюдение назначений специалиста.",

			SystolicBP:  190,
			DiastolicBP: 120,

			ImageURL: "/resources/media/stage3.png",
			VideoURL: "/resources/media/stage3.mp4",

			Status: ds.HypertensionPublished,
		},

		{
			Title: "Удаленная карточка II стадии",

			Description: "Тестовая логически удаленная услуга",

			FullDescription: "Эта запись нужна для демонстрации статуса удален и не должна отображаться в приложении.",

			SystolicBP:  170,
			DiastolicBP: 105,

			ImageURL: "/resources/media/deleted_stage2.png",
			VideoURL: "/resources/media/deleted_stage2.mp4",

			Status: ds.HypertensionDeleted,
		},
	}

	// Создаём карточки только если их ещё нет, чтобы повторный запуск seed
	// не плодил дубликаты. Сохраняем реальные ID из БД для таблицы лайков.
	servicesByTitle := make(map[string]ds.HypertensionService)

	for _, service := range services {
		var stored ds.HypertensionService
		result := db.Where("title = ?", service.Title).FirstOrCreate(&stored, service)
		if result.Error != nil {
			panic(result.Error)
		}

		servicesByTitle[stored.Title] = stored
	}

	// Пользователи для демонстрации m-m связи пользователь <-> услуга.
	usernames := []string{
		"elizabeth",
		"anna",
		"maxim",
		"dmitry",
		"sofia",
		"alexey",
	}

	usersByUsername := make(map[string]ds.User)

	for _, username := range usernames {
		var user ds.User
		result := db.Where("username = ?", username).FirstOrCreate(
			&user,
			ds.User{Username: username},
		)
		if result.Error != nil {
			panic(result.Error)
		}

		usersByUsername[username] = user
	}

	// Набор лайков. Один пользователь лайкает конкретную карточку не более одного раза.
	// Количество лайков на опубликованных карточках специально разное, чтобы это было
	// хорошо видно и в приложении, и в Adminer/pgAdmin.
	//likes := []struct {
	//	Username     string
	//	ServiceTitle string
	//}{
	//	{"elizabeth", "Оптимальное давление"},
	//	{"anna", "Оптимальное давление"},
	//	{"maxim", "Оптимальное давление"},
	//
	//	{"elizabeth", "Нормальное давление"},
	//	{"anna", "Нормальное давление"},
	//	{"sofia", "Нормальное давление"},
	//	{"alexey", "Нормальное давление"},
	//
	//	{"maxim", "Высокое нормальное давление"},
	//	{"dmitry", "Высокое нормальное давление"},
	//
	//	{"elizabeth", "I стадия гипертонии"},
	//	{"anna", "I стадия гипертонии"},
	//	{"maxim", "I стадия гипертонии"},
	//	{"dmitry", "I стадия гипертонии"},
	//	{"sofia", "I стадия гипертонии"},
	//
	//	{"anna", "II стадия гипертонии"},
	//	{"alexey", "II стадия гипертонии"},
	//	{"sofia", "II стадия гипертонии"},
	//
	//	{"elizabeth", "III стадия гипертонии"},
	//	{"dmitry", "III стадия гипертонии"},
	//	{"alexey", "III стадия гипертонии"},
	//}
	//
	//for _, item := range likes {
	//	user := usersByUsername[item.Username]
	//	service := servicesByTitle[item.ServiceTitle]
	//
	//	like := ds.HypertensionLike{
	//		UserID:                user.ID,
	//		HypertensionServiceID: service.ID,
	//	}
	//
	//	result := db.Where(
	//		"user_id = ? AND hypertension_service_id = ?",
	//		like.UserID,
	//		like.HypertensionServiceID,
	//	).FirstOrCreate(&like)
	//
	//	if result.Error != nil {
	//		panic(result.Error)
	//	}
	//}
}
