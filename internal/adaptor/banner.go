package adaptor

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/response"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerBanner struct {
	Banner usecase.BannerService
	Logger *zap.Logger
}

func NewHandlerBanner(banner usecase.BannerService, logger *zap.Logger) HandlerBanner {
	return HandlerBanner{
		Banner: banner,
		Logger: logger,
	}
}

func (h *HandlerBanner) List(ctx *gin.Context) {
	var q dto.BannerListQuery
	_ = ctx.ShouldBindQuery(&q)
	res, err := h.Banner.List(ctx.Request.Context(), q)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success", res)
}

func (h *HandlerBanner) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := h.Banner.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success", res)
}

func (h *HandlerBanner) Create(ctx *gin.Context) {
	var req dto.CreateBannerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "image is required")
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()

	url, err := utils.UploadImageToCDN(ctx.Request.Context(), f, fileHeader.Filename, "ecommerce_project")
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadGateway, err.Error())
		return
	}
	req.Image = url

	if err := h.Banner.Create(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusCreated, "created", nil)
}

func (h *HandlerBanner) Update(ctx *gin.Context) {
	var req dto.UpdateBannerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if fh, err := ctx.FormFile("image"); err == nil && fh != nil {
		f, err := fh.Open()
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer f.Close()

		url, err := utils.UploadImageToCDN(ctx.Request.Context(), f, fh.Filename, "ecommerce_project")
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusBadGateway, err.Error())
			return
		}
		req.Image = url
	}

	if err := h.Banner.Update(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "updated", nil)
}

func (h *HandlerBanner) TogglePublished(ctx *gin.Context) {
	var req dto.ToggleBannerPublishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Banner.TogglePublished(ctx.Request.Context(), req.ID, req.Published); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "publish toggled", nil)
}

func (h *HandlerBanner) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Banner.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "deleted", nil)
}
