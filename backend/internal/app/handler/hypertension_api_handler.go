package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"hypertensionstages/internal/app/auth"
	"hypertensionstages/internal/app/ds"
	"hypertensionstages/internal/app/repository"
	"hypertensionstages/internal/app/storage"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	Repo  *repository.HypertensionRepository
	Auth  *auth.Manager
	Media *storage.Client
}

func NewAPIHandler(repo *repository.HypertensionRepository, a *auth.Manager, m *storage.Client) *APIHandler {
	return &APIHandler{Repo: repo, Auth: a, Media: m}
}
func (h *APIHandler) OptionalAuth(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		if user, _, err := h.Auth.Verify(strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))); err == nil {
			c.Set("userID", user)
		}
	}
	c.Next()
}
func (h *APIHandler) RequireAuth(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer JWT required"})
		return
	}
	user, claims, err := h.Auth.Verify(strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer ")))
	if err != nil {
    fmt.Printf("JWT verification failed: %v\n", err)

    c.AbortWithStatusJSON(
        http.StatusUnauthorized,
        gin.H{"error": "invalid, expired or revoked JWT"},
    )
    return
	}
	c.Set("userID", user)
	c.Set("claims", claims)
	c.Next()
}
func (h *APIHandler) stage(s ds.HypertensionStage, userID uint) gin.H {
	count := 0
	liked := false
	for _, l := range s.Likes {
		if l.Value == 1 {
			count++
			if userID != 0 && l.UserID == userID {
				liked = true
			}
		}
	}
	image := ""
	video := ""
	if s.ImageName != "" {
		image = "/s3/hypertension-media/" + s.ImageName
	}
	if s.VideoName != "" {
		video = "/s3/hypertension-media/" + s.VideoName
	}
	return gin.H{"id": s.ID, "title": s.Title, "description": s.Description,
		"systolic_bp": s.SystolicBP, "diastolic_bp": s.DiastolicBP, "status": s.Status,
		"image_name": s.ImageName, "video_name": s.VideoName, "image_url": image, "video_url": video,
		"user_id": s.UserID, "is_owner": userID != 0 && s.UserID == userID,
		"is_liked": liked, "likes_count": count, "created_at": s.CreatedAt}
}
func (h *APIHandler) GetStages(c *gin.Context) {
	min, max := 0, 300
	if raw := c.Query("min_sbp"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid min_sbp"})
			return
		}
		min = v
	}
	if raw := c.Query("max_sbp"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid max_sbp"})
			return
		}
		max = v
	}
	if min > max {
		c.JSON(400, gin.H{"error": "min_sbp greater than max_sbp"})
		return
	}
	title := strings.TrimSpace(c.Query("title"))
	if title == "" {
		title = strings.TrimSpace(c.Query("filter"))
	}
	stages, err := h.Repo.APIGetStages(min, max, title)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	items := make([]gin.H, 0, len(stages))
	for _, s := range stages {
		items = append(items, h.stage(s, c.GetUint("userID")))
	}
	c.JSON(200, gin.H{"hypertension_stages": items})
}
func (h *APIHandler) GetFeed(c *gin.Context) {
	stages, err := h.Repo.APIGetFeed()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(stages) == 0 {
		c.JSON(404, gin.H{"error": "no published stages"})
		return
	}
	c.JSON(200, gin.H{"hypertension_stage": h.stage(stages[0], c.GetUint("userID"))})
}
func (h *APIHandler) GetFeedByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var stage ds.HypertensionStage
	if c.Query("next") == "true" {
		stage, err = h.Repo.APINextPublished(id)
	} else {
		stage, err = h.Repo.APIPublishedByID(id)
	}
	if err != nil {
		c.JSON(404, gin.H{"error": "published stage not found"})
		return
	}
	// By-ID methods preload likes in repository to supply is_liked and likes_count.
	c.JSON(200, gin.H{"hypertension_stage": h.stage(stage, c.GetUint("userID"))})
}
func (h *APIHandler) GetDraft(c *gin.Context) {
	stage, err := h.Repo.APIGetDraft(c.GetUint("userID"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if stage.ID == 0 {
		c.JSON(200, gin.H{"draft": nil})
		return
	}
	c.JSON(200, gin.H{"draft": h.stage(stage, c.GetUint("userID"))})
}
func (h *APIHandler) CreateStage(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	sbp, e1 := strconv.Atoi(c.PostForm("systolic_bp"))
	dbp, e2 := strconv.Atoi(c.PostForm("diastolic_bp"))
	if title == "" || len(title) > 160 || description == "" || e1 != nil || e2 != nil || sbp < 40 || sbp > 300 || dbp < 20 || dbp > 250 {
		c.JSON(400, gin.H{"error": "fill title, description and valid blood pressure values"})
		return
	}
	imageName := ""
	videoName := ""
	for _, entry := range []struct {
		field string
		dest  *string
		size  int64
	}{{"image", &imageName, 10 << 20}, {"video", &videoName, 80 << 20}} {
		fh, err := c.FormFile(entry.field)
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				continue
			}
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		file, err := fh.Open()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		name, err := h.Media.Upload(file, fh.Filename, entry.size)
		_ = file.Close()
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		*entry.dest = name
	}
	stage, err := h.Repo.APICreateStage(c.GetUint("userID"), title, description, sbp, dbp, imageName, videoName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"hypertension_stage": h.stage(stage, c.GetUint("userID"))})
}
func (h *APIHandler) PublishStage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	stage, err := h.Repo.APIPublishStage(id, c.GetUint("userID"))
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"hypertension_stage": h.stage(stage, c.GetUint("userID"))})
}
func (h *APIHandler) DeleteStage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	if err = h.Repo.APISoftDeleteStage(id, c.GetUint("userID")); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "deleted"})
}
func (h *APIHandler) LikeStage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Value int `json:"value"`
	}
	if err = c.ShouldBindJSON(&body); err != nil || body.Value != 0 && body.Value != 1 {
		c.JSON(400, gin.H{"error": "value must be 0 or 1"})
		return
	}
	like, err := h.Repo.APIToggleLike(id, c.GetUint("userID"), body.Value)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"like": like})
}
func (h *APIHandler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}
	body.Username = strings.TrimSpace(body.Username)
	if len(body.Username) < 3 || len(body.Username) > 60 || len(body.Password) < 6 || len(body.Password) > 72 {
		c.JSON(400, gin.H{"error": "username 3-60 chars, password 6-72 chars"})
		return
	}
	user, err := h.Repo.APIRegisterUser(body.Username, body.Password)
	if err != nil {
		c.JSON(409, gin.H{"error": "username is taken or registration failed"})
		return
	}
	c.JSON(201, gin.H{"user": gin.H{"id": user.ID, "username": user.Username}})
}
func (h *APIHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}
	user, err := h.Repo.APIAuthenticateUser(strings.TrimSpace(body.Username), body.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}
	token, err := h.Auth.Sign(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot sign token"})
		return
	}
	c.JSON(200, gin.H{"token": token, "user": gin.H{"id": user.ID, "username": user.Username}})
}
func (h *APIHandler) Logout(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(401, gin.H{"error": "missing JWT"})
		return
	}
	if err := h.Auth.Revoke(claims.(auth.Claims)); err != nil {
		c.JSON(503, gin.H{"error": "Redis unavailable; cannot revoke token"})
		return
	}
	c.JSON(200, gin.H{"status": "logged_out"})
}
