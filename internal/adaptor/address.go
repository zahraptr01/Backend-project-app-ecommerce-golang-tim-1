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

type HandlerAddress struct {
	Service usecase.AddressService
	Logger  *zap.Logger
}

func NewHandlerAddress(s usecase.AddressService, logger *zap.Logger) HandlerAddress {
	return HandlerAddress{Service: s, Logger: logger}
}

func (h *HandlerAddress) Create(ctx *gin.Context) {
	var req dto.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Service.CreateAddress(ctx.Request.Context(), req, customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusCreated, "created", res)
}

func (h *HandlerAddress) List(ctx *gin.Context) {
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Service.ListAddresses(ctx.Request.Context(), customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "list", res)
}

func (h *HandlerAddress) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	var req dto.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Service.UpdateAddress(ctx.Request.Context(), id, req, customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "updated", res)
}

func (h *HandlerAddress) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	if err := h.Service.DeleteAddress(ctx.Request.Context(), id, customerID); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "deleted", nil)
}

func (h *HandlerAddress) SetDefault(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	if err := h.Service.SetDefaultAddress(ctx.Request.Context(), customerID, id); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "set default", nil)
}
