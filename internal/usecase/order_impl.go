package usecase

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"time"
)

type orderService struct {
	repo   repository.Repository
	logger *zap.Logger
}

func NewOrderService(repo repository.Repository, logger *zap.Logger) OrderService {
	return &orderService{repo: repo, logger: logger}
}

func (s *orderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest, customerID uint) (*dto.OrderResponse, error) {
	// get cart
	cart, err := s.repo.CartRepo.GetCartByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// calculate total
	var total float64
	items := make([]entity.OrderItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		total += float64(it.Quantity) * it.UnitPrice
		items = append(items, entity.OrderItem{ProductVariantID: it.ProductVariantID, Quantity: it.Quantity, UnitPrice: it.UnitPrice})
	}

	// prepare order
	order := &entity.Order{
		CustomerID:    customerID,
		AddressID:     req.AddressID,
		Note:          "",
		PaymentMethod: req.PaymentMethod,
		Status:        "created",
		Items:         items,
	}

	// apply voucher if present
	var discount float64
	if req.VoucherCode != nil && *req.VoucherCode != "" {
		promo, err := s.repo.PromotionRepo.GetByVoucherCode(ctx, *req.VoucherCode)
		if err != nil {
			return nil, errors.New("invalid voucher")
		}
		now := time.Now()
		if !promo.Published || promo.StartDate.After(now) || promo.EndDate.Before(now) {
			return nil, errors.New("voucher not active")
		}

		// enforce usage limit
		if promo.UsageLimit <= 0 {
			return nil, errors.New("voucher usage limit exceeded")
		}
		// compute discount based on type
		if promo.Type == "percentage" {
			discount = total * promo.Discount / 100.0
		} else {
			// assume fixed amount
			discount = promo.Discount
		}
		if discount < 0 {
			discount = 0
		}
		order.VoucherCode = req.VoucherCode
		order.Discount = discount
		// store promotion id for audit
		order.PromotionID = &promo.ID
	}

	totalAfter := total - discount
	if totalAfter < 0 {
		totalAfter = 0
	}

	// decrement promotion usage if applicable (do it before creating order to avoid races)
	if order.PromotionID != nil {
		if err := s.repo.PromotionRepo.DecrementUsage(ctx, *order.PromotionID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.OrderRepo.CreateOrder(ctx, order); err != nil {
		// try to rollback usage decrement is not possible here; log and return error
		s.logger.Error("failed to create order after decrementing promotion", zap.Error(err))
		return nil, err
	}

	// clear cart
	if err := s.repo.CartRepo.ClearCart(ctx, customerID); err != nil {
		s.logger.Warn("failed to clear cart", zap.Error(err))
	}

	// build response
	respItems := make([]dto.OrderItemDTO, 0, len(order.Items))
	for _, it := range order.Items {
		respItems = append(respItems, dto.OrderItemDTO{ProductVariantID: it.ProductVariantID, Quantity: it.Quantity, UnitPrice: it.UnitPrice})
	}

	return &dto.OrderResponse{ID: order.ID, Items: respItems, Total: totalAfter, Status: order.Status}, nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, id uint, customerID uint) (*dto.OrderResponse, error) {
	o, err := s.repo.OrderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if o.CustomerID != customerID {
		return nil, errors.New("not allowed")
	}
	respItems := make([]dto.OrderItemDTO, 0, len(o.Items))
	var total float64
	for _, it := range o.Items {
		respItems = append(respItems, dto.OrderItemDTO{ProductVariantID: it.ProductVariantID, Quantity: it.Quantity, UnitPrice: it.UnitPrice})
		total += float64(it.Quantity) * it.UnitPrice
	}
	return &dto.OrderResponse{ID: o.ID, Items: respItems, Total: total, Status: o.Status}, nil
}

func (s *orderService) ListOrderHistory(ctx context.Context, customerID uint, limit, offset int) ([]dto.OrderResponse, int64, error) {
	orders, total, err := s.repo.OrderRepo.ListOrdersByCustomer(ctx, customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	res := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		var totalF float64
		items := make([]dto.OrderItemDTO, 0, len(o.Items))
		for _, it := range o.Items {
			items = append(items, dto.OrderItemDTO{ProductVariantID: it.ProductVariantID, Quantity: it.Quantity, UnitPrice: it.UnitPrice})
			totalF += float64(it.Quantity) * it.UnitPrice
		}
		res = append(res, dto.OrderResponse{ID: o.ID, Items: items, Total: totalF, Status: o.Status})
	}
	return res, total, nil
}

func (s *orderService) GetCart(ctx context.Context, customerID uint) (*dto.CartResponse, error) {
	cart, err := s.repo.CartRepo.GetCartByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	resItems := make([]dto.CartItemResponse, 0, len(cart.Items))
	var total float64
	for _, it := range cart.Items {
		resItems = append(resItems, dto.CartItemResponse{ProductVariantID: it.ProductVariantID, Quantity: it.Quantity, UnitPrice: it.UnitPrice})
		total += float64(it.Quantity) * it.UnitPrice
	}
	return &dto.CartResponse{CustomerID: cart.CustomerID, Items: resItems, Total: total}, nil
}
