package adaptor

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerCustomer struct {
	Customer usecase.CustomerService
	Logger   *zap.Logger
}

func NewHandlerCustomer(customer usecase.CustomerService, logger *zap.Logger) HandlerCustomer {
	return HandlerCustomer{
		Customer: customer,
		Logger:   logger,
	}
}

func (h *HandlerCustomer) RegisterCustomer(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.Customer.RegisterCustomer(ctx.Request.Context(), req)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.ResponseSuccess(ctx, http.StatusCreated, "registered successfully", res)
}
