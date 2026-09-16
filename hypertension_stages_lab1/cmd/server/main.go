package main

import (
	"hypertensionstages/internal/app/dsn"
	"hypertensionstages/internal/app/handler"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	_ = godotenv.Load()

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)

	if err != nil {
		panic("failed to connect database")
	}

	hypertensionRepository :=
		repository.NewHypertensionRepository(db)

	hypertensionHandler :=
		handler.NewHypertensionHandler(
			hypertensionRepository,
		)

	router := gin.Default()

	// HTML шаблоны
	router.LoadHTMLGlob("./templates/*")

	// CSS, картинки, видео
	router.Static(
		"/resources",
		"./resources",
	)

	// Лента: /feed/ и /feed/<id> одним GET-маршрутом
	router.GET(
		"/feed/*serviceID",
		hypertensionHandler.ShowHypertensionFeed,
	)

	// Черновик
	router.GET(
		"/draft",
		hypertensionHandler.ShowHypertensionDraft,
	)

	// Публикация черновика через ORM
	router.POST(
		"/draft/publish",
		hypertensionHandler.PublishHypertensionDraft,
	)

	// Плитка стадий
	router.GET(
		"/stages",
		hypertensionHandler.ShowHypertensionGrid,
	)

	// Логическое удаление карточки через чистый SQL UPDATE
	router.POST(
		"/stages/:serviceID/delete",
		hypertensionHandler.DeleteHypertensionService,
	)

	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
