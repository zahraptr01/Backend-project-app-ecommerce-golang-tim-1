package usecase

import (
	"context"
	"testing"
	"time"

	"errors"
	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
)

// Test that voucher with zero usage limit is rejected
func TestCreateOrder_VoucherUsageLimitZero(t *testing.T){
	cart := &entity.Cart{CustomerID:1, Items: []entity.CartItem{{ProductVariantID:1, Quantity:1, UnitPrice:100}}}
	promo := &entity.Promotion{Model: entity.Model{ID:2}, Type: "fixed", Discount: 10, StartDate: time.Now().Add(-time.Hour), EndDate: time.Now().Add(time.Hour), Published:true, UsageLimit:0}
	repoVal := repository.Repository{CartRepo: &simpleCartRepo{cart:cart}, PromotionRepo: &simplePromoRepo{promo:promo}, OrderRepo: &simpleOrderRepo{}, AddressRepo: nil}
	logger, _ := zap.NewDevelopment()
	svc := &orderService{repo: repoVal, logger: logger}
	code := "PROMOZERO"
	promo.VoucherCode = &code
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
	_, err := svc.CreateOrder(context.Background(), req, 1)
	if err == nil { t.Fatalf("expected error for voucher with zero usage limit") }
	if err.Error() != "voucher usage limit exceeded" { t.Fatalf("unexpected error: %v", err) }
}

// Test that DecrementUsage is called (we simulate by using a promo repo that records calls)
type trackingPromoRepo struct{
	promo *entity.Promotion
	called bool
}
func (r *trackingPromoRepo) GetByVoucherCode(ctx context.Context, code string) (*entity.Promotion, error){
	if r.promo == nil { return nil, errors.New("not found") }
	return r.promo, nil
}
func (r *trackingPromoRepo) DecrementUsage(ctx context.Context, id uint) error { r.called = true; return nil }

func TestCreateOrder_DecrementUsageCalled(t *testing.T){
	cart := &entity.Cart{CustomerID:1, Items: []entity.CartItem{{ProductVariantID:1, Quantity:1, UnitPrice:100}}}
	promo := &entity.Promotion{Model: entity.Model{ID:3}, Type: "fixed", Discount: 10, StartDate: time.Now().Add(-time.Hour), EndDate: time.Now().Add(time.Hour), Published:true, UsageLimit:5}
	repoPromo := &trackingPromoRepo{promo: promo}
	repoVal := repository.Repository{CartRepo: &simpleCartRepo{cart:cart}, PromotionRepo: repoPromo, OrderRepo: &simpleOrderRepo{}, AddressRepo: nil}
	logger, _ := zap.NewDevelopment()
	svc := &orderService{repo: repoVal, logger: logger}
	code := "PROMOTRACK"
	promo.VoucherCode = &code
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
	_, err := svc.CreateOrder(context.Background(), req, 1)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if !repoPromo.called { t.Fatalf("expected DecrementUsage to be called") }
}
