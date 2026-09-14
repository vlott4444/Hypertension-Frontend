package api

import (
	"log"

	"hypertensionstages/internal/app/handler"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartHypertensionServer() {
	log.Println("Starting hypertension stages server")

	hypertensionRepository, err := repository.NewHypertensionRepository()
	if err != nil {
		logrus.Fatal("ошибка инициализации коллекции стадий гипертонии: ", err)
	}

	hypertensionHandler := handler.NewHypertensionHandler(hypertensionRepository)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")

	router.GET("/feed/*serviceID", hypertensionHandler.ShowHypertensionFeed)
	router.GET("/draft", hypertensionHandler.ShowHypertensionDraft)
	router.GET("/stages", hypertensionHandler.ShowHypertensionGrid)

	if err := router.Run(":8080"); err != nil {
		logrus.Fatal("сервер стадий гипертонии остановлен с ошибкой: ", err)
	}
}
