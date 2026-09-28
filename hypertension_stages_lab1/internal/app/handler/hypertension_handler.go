package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"hypertensionstages/internal/app/ds"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type HypertensionHandler struct {
	HypertensionRepository *repository.HypertensionRepository
}

type HypertensionCardView struct {
	Service         ds.HypertensionService
	StageTitle      string
	StageRomanTitle string
	StageRange      string
	LikesCount      int
}

func NewHypertensionHandler(
	hypertensionRepository *repository.HypertensionRepository,
) *HypertensionHandler {
	return &HypertensionHandler{
		HypertensionRepository: hypertensionRepository,
	}
}

// GET — лента
func (h *HypertensionHandler) ShowHypertensionFeed(ctx *gin.Context) {
	service, err := h.getHypertensionServiceForFeed(ctx)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "hypertension_feed.html", gin.H{
		"Card": h.BuildHypertensionCardView(service),
	})
}

// GET — черновик
func (h *HypertensionHandler) ShowHypertensionDraft(ctx *gin.Context) {
	draft, err := h.HypertensionRepository.HypertensionDraftService()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "hypertension_draft.html", gin.H{
		"Card": h.BuildHypertensionCardView(draft),
	})
}

// POST — публикация через ORM
func (h *HypertensionHandler) PublishHypertensionDraft(ctx *gin.Context) {
	serviceID, systolicBP, diastolicBP, err := parseDraftForm(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, "некорректные данные")
		return
	}

	err = h.HypertensionRepository.PublishHypertensionDraft(
		serviceID,
		ctx.PostForm("description"),
		systolicBP,
		diastolicBP,
	)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stages")
}

// GET — плитка и поиск
func (h *HypertensionHandler) ShowHypertensionGrid(ctx *gin.Context) {
	services, boundsMin, boundsMax, sbpMin, sbpMax, filterError, err :=
		h.getHypertensionServicesForGrid(ctx)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	cards := make([]HypertensionCardView, 0, len(services))

	for _, service := range services {
		cards = append(cards, h.BuildHypertensionCardView(service))
	}

	ctx.HTML(http.StatusOK, "hypertension_grid.html", gin.H{
		"Cards": cards,
		"SBPBounds": gin.H{
			"Min": boundsMin,
			"Max": boundsMax,
		},
		"SBPMin":      sbpMin,
		"SBPMax":      sbpMax,
		"SBPQuery":    fmt.Sprintf("%d–%d", sbpMin, sbpMax),
		"FilterError": filterError,
	})
}

// POST — удаление чистый SQL
func (h *HypertensionHandler) DeleteHypertensionService(ctx *gin.Context) {
	serviceID, err := strconv.Atoi(ctx.Param("serviceID"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "некорректный id")
		return
	}

	err = h.HypertensionRepository.DeleteHypertensionServiceSQL(
		ctx.Request.Context(),
		serviceID,
	)

	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/stages")
}

// --------------------
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// --------------------

func (h *HypertensionHandler) getHypertensionServiceForFeed(
	ctx *gin.Context,
) (ds.HypertensionService, error) {

	rawID := strings.Trim(ctx.Param("serviceID"), "/")

	if rawID == "" {
		services, err :=
			h.HypertensionRepository.PublishedHypertensionServices()

		if err != nil {
			return ds.HypertensionService{}, err
		}

		return services[0], nil
	}

	serviceID, err := strconv.Atoi(rawID)
	if err != nil {
		return ds.HypertensionService{}, err
	}

	if ctx.Query("next") == "true" {
		return h.HypertensionRepository.
			NextPublishedHypertensionService(serviceID)
	}

	return h.HypertensionRepository.
		HypertensionServiceByID(serviceID)
}

func parseDraftForm(
	ctx *gin.Context,
) (int, int, int, error) {

	serviceID, err := strconv.Atoi(ctx.PostForm("service_id"))
	if err != nil {
		return 0, 0, 0, err
	}

	systolicBP, err := strconv.Atoi(ctx.PostForm("sbp"))
	if err != nil {
		return 0, 0, 0, err
	}

	diastolicBP, err := strconv.Atoi(ctx.PostForm("dbp"))
	if err != nil {
		return 0, 0, 0, err
	}

	return serviceID, systolicBP, diastolicBP, nil
}

// getHypertensionServicesForGrid возвращает карточки, отфильтрованные
// по диапазону [sbp_min, sbp_max], а также границы шкалы и текущие
// значения ползунков для шаблона.
func (h *HypertensionHandler) getHypertensionServicesForGrid(
	ctx *gin.Context,
) (
	services []ds.HypertensionService,
	boundsMin, boundsMax int,
	sbpMin, sbpMax int,
	filterError string,
	err error,
) {
	// 1. Границы шкалы — минимум и максимум САД по всем
	//    опубликованным карточкам.
	boundsMin, boundsMax, err = h.HypertensionRepository.SBPBounds()
	if err != nil {
		return nil, 0, 0, 0, 0, "", err
	}

	if boundsMin == 0 && boundsMax == 0 {
		boundsMin, boundsMax = 100, 200
	}

	// 2. Текущие значения ползунков. По умолчанию — границы шкалы.
	sbpMin = boundsMin
	sbpMax = boundsMax

	if raw := strings.TrimSpace(ctx.Query("sbp_min")); raw != "" {
		v, convErr := strconv.Atoi(raw)
		if convErr != nil {
			return nil, boundsMin, boundsMax, boundsMin, boundsMax,
				"Минимальное САД должно быть целым числом", nil
		}
		sbpMin = v
	}

	if raw := strings.TrimSpace(ctx.Query("sbp_max")); raw != "" {
		v, convErr := strconv.Atoi(raw)
		if convErr != nil {
			return nil, boundsMin, boundsMax, boundsMin, boundsMax,
				"Максимальное САД должно быть целым числом", nil
		}
		sbpMax = v
	}

	// 3. Без JS ползунки могут перескочить друг через друга —
	//    просто меняем местами, чтобы фильтр остался корректным.
	if sbpMin > sbpMax {
		sbpMin, sbpMax = sbpMax, sbpMin
	}

	// 4. Фильтрация по диапазону.
	services, err = h.HypertensionRepository.
		FilterPublishedHypertensionBySBPRange(sbpMin, sbpMax)

	return services, boundsMin, boundsMax, sbpMin, sbpMax, "", err
}

func (h *HypertensionHandler) BuildHypertensionCardView(
	service ds.HypertensionService,
) HypertensionCardView {

	return HypertensionCardView{
		Service: service,

		StageTitle: ds.HypertensionStageTitleBySBP(
			service.SystolicBP,
		),

		StageRomanTitle: ds.HypertensionStageRomanTitleBySBP(
			service.SystolicBP,
		),

		StageRange: ds.HypertensionStageRangeBySBP(
			service.SystolicBP,
		),

		LikesCount: len(service.Likes),
	}
}
