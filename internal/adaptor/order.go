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

type HandlerOrder struct {
	Order  usecase.OrderService
	Logger *zap.Logger
}

func NewHandlerOrder(order usecase.OrderService, logger *zap.Logger) HandlerOrder {
	return HandlerOrder{Order: order, Logger: logger}
}

func (h *HandlerOrder) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Order.CreateOrder(ctx.Request.Context(), req, customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusCreated, "created", res)
}

func (h *HandlerOrder) GetOrderDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Order.GetOrderDetail(ctx.Request.Context(), id, customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "detail", res)
}

func (h *HandlerOrder) ListOrderHistory(ctx *gin.Context) {
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	page := 1
	limit := 20
	if p := ctx.Query("page"); p != "" {
		if pi, err := strconv.Atoi(p); err == nil {
			page = pi
		}
	}
	offset := (page - 1) * limit
	res, total, err := h.Order.ListOrderHistory(ctx.Request.Context(), customerID, limit, offset)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "history", gin.H{"data": res, "total": total})
}

// Cart returns current cart summary (simple proxy to history for now)
func (h *HandlerOrder) Cart(ctx *gin.Context) {
	uid, _ := ctx.Get("userID")
	customerID, _ := uid.(uint)
	res, err := h.Order.GetCart(ctx.Request.Context(), customerID)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "cart", res)
}
