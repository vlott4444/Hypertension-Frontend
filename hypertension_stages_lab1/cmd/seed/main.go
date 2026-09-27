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

	// -------------------------------------------------------
	// 1. ПОЛЬЗОВАТЕЛИ
	// -------------------------------------------------------

	users := []ds.User{
		{
			Username: "elizabeth",
			Password: "pass1",
		},
		{
			Username: "anna",
			Password: "pass2",
		},
		{
			Username: "maxim",
			Password: "pass3",
		},
		{
			Username: "dmitry",
			Password: "pass4",
		},
		{
			Username: "sofia",
			Password: "pass5",
		},
		{
			Username: "alexey",
			Password: "pass6",
		},
	}

	usersByUsername := make(map[string]ds.User)

	for _, seedUser := range users {
		var user ds.User

		// Если пользователь уже существует,
		// Assign обновит ему пароль.
		result := db.
			Where("username = ?", seedUser.Username).
			Assign(ds.User{
				Password: seedUser.Password,
			}).
			FirstOrCreate(&user, seedUser)

		if result.Error != nil {
			panic(result.Error)
		}

		usersByUsername[user.Username] = user
	}

	// Все тестовые карточки принадлежат Elizabeth.
	owner := usersByUsername["elizabeth"]

	// -------------------------------------------------------
	// 2. MINIO
	// -------------------------------------------------------

	const minioBaseURL = "http://localhost:9000/hypertension-media"

	// -------------------------------------------------------
	// 3. КАРТОЧКИ
	// -------------------------------------------------------

	services := []ds.HypertensionService{
		{
			Title: "Оптимальное давление",

			Description: "Систолическое давление ниже 120 мм рт. ст.",

			FullDescription: "Оптимальный уровень артериального давления характеризуется значениями ниже 120 мм рт. ст. для систолического давления и ниже 80 мм рт. ст. для диастолического. Обычно такое давление связано с низким риском сердечно-сосудистых осложнений. Для поддержания нормальных показателей рекомендуется сохранять физическую активность, придерживаться сбалансированного питания и регулярно контролировать давление.",

			SystolicBP:  110,
			DiastolicBP: 70,

			ImageURL: minioBaseURL + "/stage_optimal.png",
			VideoURL: minioBaseURL + "/optimal.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		{
			Title: "Нормальное давление",

			Description: "Систолическое давление 120–129 мм рт. ст.",

			FullDescription: "Нормальное артериальное давление находится в диапазоне 120–129 мм рт. ст. по систолическому показателю. Такие значения считаются благоприятными, однако важно следить за динамикой давления. Регулярные измерения помогают вовремя заметить изменения и предотвратить развитие стойкой гипертонии.",

			SystolicBP:  125,
			DiastolicBP: 78,

			ImageURL: minioBaseURL + "/stage_normal.png",
			VideoURL: minioBaseURL + "/normal.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		{
			Title: "Высокое нормальное давление",

			Description: "Систолическое давление 130–139 мм рт. ст.",

			FullDescription: "Повышенное нормальное давление означает, что показатели уже находятся выше оптимального уровня. Это состояние не всегда требует лечения, но является сигналом для изменения образа жизни. Важно контролировать массу тела, уровень физической активности, питание и регулярно измерять давление.",

			SystolicBP:  135,
			DiastolicBP: 85,

			ImageURL: minioBaseURL + "/stage_high_normal.png",
			VideoURL: minioBaseURL + "/high_normal.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		{
			Title: "I стадия гипертонии",

			Description: "Систолическое давление 140–159 мм рт. ст.",

			FullDescription: "Первая стадия артериальной гипертонии характеризуется устойчивым повышением давления. На этом этапе могут отсутствовать выраженные симптомы, поэтому регулярный контроль особенно важен. Врач оценивает общий риск осложнений и при необходимости назначает лечение вместе с рекомендациями по изменению образа жизни.",

			SystolicBP:  150,
			DiastolicBP: 95,

			ImageURL: minioBaseURL + "/stage1.png",
			VideoURL: minioBaseURL + "/stage1.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		{
			Title: "II стадия гипертонии",

			Description: "Систолическое давление 160–179 мм рт. ст.",

			FullDescription: "Вторая стадия гипертонии характеризуется более значительным повышением давления. При длительном сохранении таких значений увеличивается нагрузка на сердце, сосуды и другие органы. Обычно требуется регулярное наблюдение специалиста и подбор терапии для достижения целевых показателей давления.",

			SystolicBP:  170,
			DiastolicBP: 105,

			ImageURL: minioBaseURL + "/stage2.png",
			VideoURL: minioBaseURL + "/stage2.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		{
			Title: "III стадия гипертонии",

			Description: "Систолическое давление 180 мм рт. ст. и выше",

			FullDescription: "Третья стадия артериальной гипертонии характеризуется выраженным повышением давления. Такие значения требуют особого внимания, поскольку связаны с повышенной нагрузкой на сердечно-сосудистую систему. Необходимо регулярное медицинское наблюдение и соблюдение назначений специалиста.",

			SystolicBP:  190,
			DiastolicBP: 120,

			ImageURL: minioBaseURL + "/stage3.png",
			VideoURL: minioBaseURL + "/stage3.mp4",

			Status: ds.HypertensionPublished,
			UserID: owner.ID,
		},

		// -------------------------------------------------------
		// ЧЕРНОВИК
		// -------------------------------------------------------

		{
			Title: "Черновик I стадии",

			Description: "Тестовая карточка в статусе черновик",

			FullDescription: "Эта запись используется для демонстрации карточки в статусе черновик.",

			SystolicBP:  150,
			DiastolicBP: 95,

			ImageURL: minioBaseURL + "/draft_stage1.png",
			VideoURL: minioBaseURL + "/draft_stage1.mp4",

			Status: ds.HypertensionDraft,
			UserID: owner.ID,
		},

		// -------------------------------------------------------
		// ЛОГИЧЕСКИ УДАЛЁННАЯ КАРТОЧКА
		// -------------------------------------------------------

		{
			Title: "Удаленная карточка II стадии",

			Description: "Тестовая логически удаленная услуга",

			FullDescription: "Эта запись нужна для демонстрации статуса удален и не должна отображаться в приложении.",

			SystolicBP:  170,
			DiastolicBP: 105,

			ImageURL: minioBaseURL + "/deleted_stage2.png",
			VideoURL: minioBaseURL + "/deleted_stage2.mp4",

			Status: ds.HypertensionDeleted,
			UserID: owner.ID,
		},
	}

	// -------------------------------------------------------
	// 4. СОЗДАЁМ ИЛИ ОБНОВЛЯЕМ КАРТОЧКИ
	// -------------------------------------------------------

	for _, service := range services {
		var stored ds.HypertensionService

		result := db.
			Where("title = ?", service.Title).
			Assign(service).
			FirstOrCreate(&stored)

		if result.Error != nil {
			panic(result.Error)
		}
	}

	// -------------------------------------------------------
	// 5. ЛАЙКИ
	// -------------------------------------------------------
	//
	// Во второй лабораторной ставить лайки через приложение
	// не требуется.
	//
	// Если понадобятся тестовые записи m-m,
	// этот блок можно добавить позже.
}
