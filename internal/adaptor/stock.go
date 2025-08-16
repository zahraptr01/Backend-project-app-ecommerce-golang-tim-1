package adaptor

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerStock struct {
	Stock  usecase.StockService
	Logger *zap.Logger
}

func NewHandlerStock(stock usecase.StockService, logger *zap.Logger) HandlerStock {
	return HandlerStock{
		Stock:  stock,
		Logger: logger,
	}
}

func (h *HandlerStock) List(ctx *gin.Context) {
	var q dto.StockListQuery
	_ = ctx.ShouldBindQuery(&q)

	res, err := h.Stock.List(ctx.Request.Context(), q)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "ok", res)
}

func (h *HandlerStock) Detail(ctx *gin.Context) {
	idStr := ctx.Param("variant_id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "invalid variant_id")
		return
	}
	res, err := h.Stock.Detail(ctx.Request.Context(), uint(id))
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "ok", res)
}

func (h *HandlerStock) Add(ctx *gin.Context) {
	var req dto.AddStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Stock.Add(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "stock increased", nil)
}

func (h *HandlerStock) Set(ctx *gin.Context) {
	var req dto.SetStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Stock.Set(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "stock set", nil)
}

func (h *HandlerStock) Delete(ctx *gin.Context) {
	var req dto.DeleteStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Stock.Delete(ctx.Request.Context(), req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "stock deleted (set 0)", nil)
}

func (h *HandlerStock) VariantsDropdown(ctx *gin.Context) {
	var q dto.VariantDropdownQuery
	_ = ctx.ShouldBindQuery(&q)

	res, err := h.Stock.VariantsDropdown(ctx.Request.Context(), q)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "ok", res)
}
