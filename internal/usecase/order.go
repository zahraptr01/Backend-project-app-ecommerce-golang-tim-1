package usecase

import (
	"context"
	"project-app-ecommerce-golang-tim-1/internal/dto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest, customerID uint) (*dto.OrderResponse, error)
	GetOrderDetail(ctx context.Context, id uint, customerID uint) (*dto.OrderResponse, error)
	ListOrderHistory(ctx context.Context, customerID uint, limit, offset int) ([]dto.OrderResponse, int64, error)
	GetCart(ctx context.Context, customerID uint) (*dto.CartResponse, error)
}

// Implementation will be added later; this is a placeholder interface to wire handlers.
