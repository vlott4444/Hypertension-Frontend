package main

import (
	"hypertensionstages/internal/app/auth"
	"hypertensionstages/internal/app/dsn"
	"hypertensionstages/internal/app/handler"
	"hypertensionstages/internal/app/repository"
	"hypertensionstages/internal/app/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	hypertensionRepository := repository.NewHypertensionRepository(db)

	hypertensionHandler := handler.NewHypertensionHandler(hypertensionRepository)

	manager, err := auth.New()
	if err != nil {
		panic(err)
	}
	apiHandler := handler.NewAPIHandler(hypertensionRepository, manager, storage.New())

	router := gin.Default()

	router.LoadHTMLGlob("./templates/*")
	router.Static("/resources", "./resources")

	router.GET("/feed/*stageID", hypertensionHandler.ShowHypertensionFeed)
	// old SSR draft page removed: requires JWT and is available at the React /add route
	// Legacy SSR write routes removed: they bypassed JWT verification.
	router.GET("/stages", hypertensionHandler.ShowHypertensionGrid)

	api := router.Group("/api")
	{

		api.GET("/hypertension_stages", apiHandler.OptionalAuth, apiHandler.GetStages)
		api.GET("/feed", apiHandler.OptionalAuth, apiHandler.GetFeed)
		api.GET("/feed/:id", apiHandler.OptionalAuth, apiHandler.GetFeedByID)
		api.GET("/draft", apiHandler.RequireAuth, apiHandler.GetDraft)
		api.POST("/hypertension_stages", apiHandler.RequireAuth, apiHandler.CreateStage)
		api.PUT("/hypertension_stages/:id/publish", apiHandler.RequireAuth, apiHandler.PublishStage)
		api.DELETE("/hypertension_stages/:id", apiHandler.RequireAuth, apiHandler.DeleteStage)
		api.POST("/hypertension_stages/:id/like", apiHandler.RequireAuth, apiHandler.LikeStage)

		api.POST("/auth/register", apiHandler.Register)
		api.POST("/auth/login", apiHandler.Login)
		api.POST("/auth/logout", apiHandler.RequireAuth, apiHandler.Logout)
	}

	router.GET("/openapi.json", func(c *gin.Context) {
    c.File("swagger/doc.json")
	})
	router.GET(
    	"/swagger/*any",
    	ginSwagger.WrapHandler(
        	swaggerFiles.Handler,
        	ginSwagger.URL("/openapi.json"),
    	),
	)

	router.Run(":8080")
}
