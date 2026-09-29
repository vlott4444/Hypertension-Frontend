package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hypertensionstages/internal/app/ds"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
)

// APIHandler — REST-хендлер для /api/*.
// Всё общение — только JSON, без HTML-шаблонов.
type APIHandler struct {
	Repo *repository.HypertensionRepository
}

func NewAPIHandler(repo *repository.HypertensionRepository) *APIHandler {
	return &APIHandler{Repo: repo}
}

// ============================================================
// ДОМЕН УСЛУГИ
// ============================================================

// GET /api/services?filter=...
// Список опубликованных с фильтрацией.
// Поле "is_mine" 0/1 — если создатель совпадает с текущим пользователем.
func (h *APIHandler) GetServices(c *gin.Context) {
	user := ds.GetCurrentUser()
	filter := c.Query("filter")

	services, err := h.Repo.APIGetServices(user.ID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

// GET /api/feed
// Лента — только опубликованные.
func (h *APIHandler) GetFeed(c *gin.Context) {
	services, err := h.Repo.APIGetFeed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"feed": services})
}

// GET /api/feed/:id?next=true
// Публикация по ID. Если next=true — следующая по ленте.
func (h *APIHandler) GetFeedByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	next := c.Query("next") == "true"

	var service ds.HypertensionService
	if next {
		service, err = h.Repo.APINextPublished(id)
	} else {
		service, err = h.Repo.APIPublishedByID(id)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// GET /api/draft
// Черновик текущего пользователя. Не более одного. ID не указывается.
func (h *APIHandler) GetDraft(c *gin.Context) {
	user := ds.GetCurrentUser()

	service, err := h.Repo.APIGetDraft(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": service})
}

// POST /api/services
// Создание услуги с загрузкой файлов (multipart/form-data).
// Название услуги, описание, САД, ДАД — из формы.
// Файлы: image, video.
// Системные поля (ID, статус, создатель, даты) с клиента игнорируются.
func (h *APIHandler) CreateService(c *gin.Context) {
	user := ds.GetCurrentUser()

	title := c.PostForm("title")
	description := c.PostForm("description")
	sbp, _ := strconv.Atoi(c.PostForm("systolic_bp"))
	dbp, _ := strconv.Atoi(c.PostForm("diastolic_bp"))

	imageName, err := saveUploadedFile(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoName, err := saveUploadedFile(c, "video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.Repo.APICreateService(
		user.ID, title, description, sbp, dbp, imageName, videoName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service": service})
}

// PUT /api/services/:id/publish
// Публикация черновика. Только свои услуги.
func (h *APIHandler) PublishService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user := ds.GetCurrentUser()

	service, err := h.Repo.APIPublishService(id, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// DELETE /api/services/:id
// Soft delete. Только свои услуги.
func (h *APIHandler) DeleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user := ds.GetCurrentUser()

	if err := h.Repo.APISoftDeleteService(id, user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// POST /api/services/:id/like
// Тело: {"value": 0|1}
func (h *APIHandler) LikeService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Value int `json:"value"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if body.Value != 0 && body.Value != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value must be 0 or 1"})
		return
	}

	user := ds.GetCurrentUser()

	like, err := h.Repo.APIToggleLike(id, user.ID, body.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"like": like})
}

// ============================================================
// ДОМЕН ПОЛЬЗОВАТЕЛЯ
// ============================================================

// POST /api/auth/register
func (h *APIHandler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	user, err := h.Repo.APIRegisterUser(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// POST /api/auth/login — заглушка
func (h *APIHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "stub"})
}

// POST /api/auth/logout — заглушка
func (h *APIHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "stub"})
}

// saveUploadedFile принимает файл из формы, генерирует латинское имя
// и загружает в Minio. Возвращает имя файла (для сохранения в БД).
func saveUploadedFile(c *gin.Context, field string) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("файл %s не передан", field)
	}

	// TODO: здесь подключишь Minio-клиент из resources/minio
	// и загрузишь файл. Пока просто возвращаем имя.
	return file.Filename, nil
}
