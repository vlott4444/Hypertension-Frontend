package api

// Compatibility entry point for older labs. The main entry point is cmd/server.
import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hypertensionstages/internal/app/auth"
	"hypertensionstages/internal/app/handler"
	"hypertensionstages/internal/app/repository"
	"hypertensionstages/internal/app/storage"
	"log"
)

func StartHypertensionServer(db *gorm.DB) {
	repo := repository.NewHypertensionRepository(db)
	a, err := auth.New()
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewAPIHandler(repo, a, storage.New())
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/resources", "./resources")
	// Read-only legacy SSR routes; all writes require API JWT.
	legacy := handler.NewHypertensionHandler(repo)
	router.GET("/feed/*stageID", legacy.ShowHypertensionFeed)
	router.GET("/stages", legacy.ShowHypertensionGrid)
	api := router.Group("/api")
	api.GET("/hypertension_stages", h.OptionalAuth, h.GetStages)
	api.GET("/feed", h.OptionalAuth, h.GetFeed)
	api.GET("/feed/:id", h.OptionalAuth, h.GetFeedByID)
	api.GET("/draft", h.RequireAuth, h.GetDraft)
	api.POST("/hypertension_stages", h.RequireAuth, h.CreateStage)
	api.PUT("/hypertension_stages/:id/publish", h.RequireAuth, h.PublishStage)
	api.DELETE("/hypertension_stages/:id", h.RequireAuth, h.DeleteStage)
	api.POST("/hypertension_stages/:id/like", h.RequireAuth, h.LikeStage)
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)
	api.POST("/auth/logout", h.RequireAuth, h.Logout)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
