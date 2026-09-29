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

	// Репозиторий
	hypertensionRepository := repository.NewHypertensionRepository(db)

	// HTML-хендлер (лабы 1-3)
	hypertensionHandler := handler.NewHypertensionHandler(hypertensionRepository)

	// API-хендлер (лаба 4)
	apiHandler := handler.NewAPIHandler(hypertensionRepository)

	router := gin.Default()

	// HTML шаблоны
	router.LoadHTMLGlob("./templates/*")
	router.Static("/resources", "./resources")

	// ---------- HTML-роуты (не трогаем) ----------
	router.GET("/feed/*serviceID", hypertensionHandler.ShowHypertensionFeed)
	router.GET("/draft", hypertensionHandler.ShowHypertensionDraft)
	router.POST("/draft/publish", hypertensionHandler.PublishHypertensionDraft)
	router.GET("/stages", hypertensionHandler.ShowHypertensionGrid)
	router.POST("/stages/:serviceID/delete", hypertensionHandler.DeleteHypertensionService)

	// ---------- REST API (лаба 4) ----------
	api := router.Group("/api")
	{
		// Домен услуги
		api.GET("/services", apiHandler.GetServices)
		api.GET("/feed", apiHandler.GetFeed)
		api.GET("/feed/:id", apiHandler.GetFeedByID)
		api.GET("/draft", apiHandler.GetDraft)
		api.POST("/services", apiHandler.CreateService)
		api.PUT("/services/:id/publish", apiHandler.PublishService)
		api.DELETE("/services/:id", apiHandler.DeleteService)
		api.POST("/services/:id/like", apiHandler.LikeService)

		// Домен пользователя
		api.POST("/auth/register", apiHandler.Register)
		api.POST("/auth/login", apiHandler.Login)
		api.POST("/auth/logout", apiHandler.Logout)
	}

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
