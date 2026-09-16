package handler

import (
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
	services, filterError, err := h.getHypertensionServicesForGrid(ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	cards := make([]HypertensionCardView, 0, len(services))

	for _, service := range services {
		cards = append(cards, h.BuildHypertensionCardView(service))
	}

	ctx.HTML(http.StatusOK, "hypertension_grid.html", gin.H{
		"Cards":       cards,
		"SBPQuery":    ctx.Query("sbp"),
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

func (h *HypertensionHandler) getHypertensionServicesForGrid(
	ctx *gin.Context,
) ([]ds.HypertensionService, string, error) {

	sbpQuery := strings.TrimSpace(ctx.Query("sbp"))

	if sbpQuery == "" {
		services, err :=
			h.HypertensionRepository.PublishedHypertensionServices()

		return services, "", err
	}

	sbp, err := strconv.Atoi(sbpQuery)
	if err != nil {
		return []ds.HypertensionService{},
			"Введите САД целым числом",
			nil
	}

	services, err :=
		h.HypertensionRepository.FilterPublishedHypertensionBySBP(sbp)

	return services, "", err
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
