package api

import (
	"log"

	"hypertensionstages/internal/app/handler"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// StartHypertensionServer оставлен как отдельная точка сборки HTTP-слоя.
// В cmd/server/main.go используется та же конфигурация маршрутов.
func StartHypertensionServer(db *gorm.DB) {
	log.Println("Starting hypertension stages server")

	hypertensionRepository := repository.NewHypertensionRepository(db)
	hypertensionHandler := handler.NewHypertensionHandler(hypertensionRepository)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/resources", "./resources")

	router.GET("/feed/*serviceID", hypertensionHandler.ShowHypertensionFeed)
	router.GET("/draft", hypertensionHandler.ShowHypertensionDraft)
	router.GET("/stages", hypertensionHandler.ShowHypertensionGrid)

	router.POST("/draft/publish", hypertensionHandler.PublishHypertensionDraft)
	router.POST("/stages/:serviceID/delete", hypertensionHandler.DeleteHypertensionService)

	if err := router.Run(":8080"); err != nil {
		logrus.Fatal("сервер стадий гипертонии остановлен с ошибкой: ", err)
	}
}
