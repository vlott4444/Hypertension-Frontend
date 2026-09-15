package handler

import (
	"net/http"
	"strconv"
	"strings"

	"hypertensionstages/internal/app/ds"
	"hypertensionstages/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (h *HypertensionHandler) ShowHypertensionFeed(ctx *gin.Context) {

	published, err := h.HypertensionRepository.PublishedHypertensionServices()

	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	rawID := strings.Trim(ctx.Param("serviceID"), "/")

	var hypertensionService ds.HypertensionService

	if rawID == "" {

		hypertensionService = published[0]

	} else {

		serviceID, parseErr := strconv.Atoi(rawID)

		if parseErr != nil {
			ctx.String(
				http.StatusBadRequest,
				"некорректный id карточки гипертонии",
			)
			return
		}

		if ctx.Query("next") == "true" {

			hypertensionService, err =
				h.HypertensionRepository.NextPublishedHypertensionService(serviceID)

		} else {

			hypertensionService, err =
				h.HypertensionRepository.HypertensionServiceByID(serviceID)

		}

		if err != nil {

			logrus.Warn(err)

			ctx.String(
				http.StatusNotFound,
				err.Error(),
			)

			return
		}
	}

	ctx.HTML(
		http.StatusOK,
		"hypertension_feed.html",
		gin.H{
			"Card": h.BuildHypertensionCardView(hypertensionService),
		},
	)
}

func (h *HypertensionHandler) ShowHypertensionDraft(
	ctx *gin.Context,
) {

	hypertensionDraft, err :=
		h.HypertensionRepository.HypertensionDraftService()

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"hypertension_draft.html",
		gin.H{
			"Card": h.BuildHypertensionCardView(hypertensionDraft),
		},
	)
}

func (h *HypertensionHandler) ShowHypertensionGrid(
	ctx *gin.Context,
) {

	sbpQuery := strings.TrimSpace(ctx.Query("sbp"))

	var hypertensionServices []ds.HypertensionService

	var err error

	filterError := ""

	if sbpQuery == "" {

		hypertensionServices, err =
			h.HypertensionRepository.PublishedHypertensionServices()

	} else {

		sbp, parseErr := strconv.Atoi(sbpQuery)

		if parseErr != nil {

			filterError =
				"Введите САД целым числом, например 150"

			hypertensionServices =
				[]ds.HypertensionService{}

		} else {

			hypertensionServices, err =
				h.HypertensionRepository.FilterPublishedHypertensionBySBP(sbp)

		}
	}

	if err != nil {

		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	cards := make(
		[]HypertensionCardView,
		0,
		len(hypertensionServices),
	)

	for _, hypertensionService := range hypertensionServices {

		cards = append(
			cards,
			h.BuildHypertensionCardView(hypertensionService),
		)

	}

	ctx.HTML(
		http.StatusOK,
		"hypertension_grid.html",
		gin.H{
			"Cards":       cards,
			"SBPQuery":    sbpQuery,
			"FilterError": filterError,
		},
	)
}

func (h *HypertensionHandler) BuildHypertensionCardView(
	hypertensionService ds.HypertensionService,
) HypertensionCardView {

	return HypertensionCardView{

		Service: hypertensionService,

		StageTitle: ds.HypertensionStageTitleBySBP(
			hypertensionService.SystolicBP,
		),

		StageRomanTitle: ds.HypertensionStageRomanTitleBySBP(
			hypertensionService.SystolicBP,
		),

		StageRange: ds.HypertensionStageRangeBySBP(
			hypertensionService.SystolicBP,
		),

		LikesCount: len(hypertensionService.Likes),
	}
}
