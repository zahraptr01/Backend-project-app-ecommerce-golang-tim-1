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

type HandlerCategory struct {
	Category usecase.CategoryService
	Logger   *zap.Logger
}

func NewHandlerCategory(category usecase.CategoryService, logger *zap.Logger) HandlerCategory {
	return HandlerCategory{
		Category: category,
		Logger:   logger,
	}
}

func (h *HandlerCategory) List(ctx *gin.Context) {
	var q dto.CategoryListQuery
	_ = ctx.ShouldBindQuery(&q)

	res, err := h.Category.List(ctx.Request.Context(), q)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success", res)
}

func (h *HandlerCategory) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	res, err := h.Category.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success", res)
}

func (h *HandlerCategory) Create(ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("image")
	name := ctx.PostForm("name")
	if name == "" || fileHeader == nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "name and image are required")
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()

	iconURL, err := utils.UploadImageToCDN(ctx.Request.Context(), f, fileHeader.Filename, "ecommerce_project")
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadGateway, err.Error())
		return
	}

	req := dto.CreateCategoryRequest{
		Name: name,
		Icon: iconURL,
	}
	if err := h.Category.Create(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusCreated, "created", gin.H{"icon": iconURL})
}

func (h *HandlerCategory) Update(ctx *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "invalid form")
		return
	}
	idStr := ctx.PostForm("id")
	name := ctx.PostForm("name")
	if idStr == "" || name == "" {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "id and name are required")
		return
	}
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	req.ID = uint(id64)
	req.Name = name

	fileHeader, _ := ctx.FormFile("image")
	if fileHeader != nil {
		f, err := fileHeader.Open()
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer f.Close()

		iconURL, err := utils.UploadImageToCDN(ctx.Request.Context(), f, fileHeader.Filename, "ecommerce_project")
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusBadGateway, err.Error())
			return
		}
		req.Icon = iconURL
	} else {
		req.Icon = ctx.PostForm("icon") // opsional; atau bisa ambil dari DB jika mau tetap
	}

	if err := h.Category.Update(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "updated", nil)
}

func (h *HandlerCategory) TogglePublished(ctx *gin.Context) {
	var req dto.TogglePublishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Category.TogglePublished(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "publish toggled", nil)
}

func (h *HandlerCategory) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Category.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "deleted", nil)
}
