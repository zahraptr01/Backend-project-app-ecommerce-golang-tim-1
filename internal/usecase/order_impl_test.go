package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
)

// Minimal mock repository to test order voucher logic
type mockRepo struct{
	repository.Repository
	Cart *entity.Cart
	Promo *entity.Promotion
}

func (m *mockRepo) CartRepo() repository.CartRepository { return nil }

func (m *mockRepo) OrderRepo() repository.OrderRepository { return nil }

func (m *mockRepo) PromotionRepo() repository.PromotionRepository { return nil }

// We will create a simplified test by creating a fake service with inline repo methods via embedding not used; instead call the logic by constructing the orderService with repo interface that has required methods.

// For succinctness, we write table-driven tests by mocking PromotionRepo and CartRepo via small structs implementing the needed methods.

// Mock CartRepo
type simpleCartRepo struct{
	cart *entity.Cart
}
func (r *simpleCartRepo) GetCartByCustomer(ctx context.Context, customerID uint) (*entity.Cart, error){
	if r.cart == nil { return nil, errors.New("no cart") }
	return r.cart, nil
}
func (r *simpleCartRepo) ClearCart(ctx context.Context, customerID uint) error { return nil }

// Mock PromotionRepo
type simplePromoRepo struct{
	promo *entity.Promotion
}
func (r *simplePromoRepo) GetByVoucherCode(ctx context.Context, code string) (*entity.Promotion, error){
	if r.promo == nil { return nil, errors.New("not found") }
	return r.promo, nil
}
func (r *simplePromoRepo) DecrementUsage(ctx context.Context, id uint) error { return nil }

// Mock OrderRepo
type simpleOrderRepo struct{}
func (r *simpleOrderRepo) CreateOrder(ctx context.Context, order *entity.Order) error { order.ID = 1; return nil }
func (r *simpleOrderRepo) GetOrderByID(ctx context.Context, id uint) (*entity.Order, error) { return nil, errors.New("not implemented") }
func (r *simpleOrderRepo) ListOrdersByCustomer(ctx context.Context, customerID uint, limit, offset int) ([]entity.Order, int64, error) { return nil, 0, errors.New("not implemented") }

// Combined repository for usecase.Repository expectation
type combinedRepo struct{
	Cart repository.CartRepository
	Promotion repository.PromotionRepository
	Order repository.OrderRepository
	Address repository.AddressRepository
}

func (r combinedRepo) RedisRepo() repository.RedisRepository { return nil }
func (r combinedRepo) AuthRepo() repository.AuthRepository { return nil }
func (r combinedRepo) CustomerRepo() repository.CustomerRepository { return nil }
func (r combinedRepo) StockRepo() repository.StockRepository { return nil }
func (r combinedRepo) CategoryRepo() repository.CategoryRepository { return nil }
func (r combinedRepo) BannerRepo() repository.BannerRepository { return nil }
func (r combinedRepo) OrderRepo() repository.OrderRepository { return r.Order }
func (r combinedRepo) AddressRepo() repository.AddressRepository { return r.Address }
func (r combinedRepo) CartRepo() repository.CartRepository { return r.Cart }
func (r combinedRepo) PromotionRepo() repository.PromotionRepository { return r.Promotion }
func (r combinedRepo) GetDB() interface{} { return nil }

func TestCreateOrder_WithVoucherPercentage(t *testing.T){
	cart := &entity.Cart{CustomerID:1, Items: []entity.CartItem{{ProductVariantID:1, Quantity:2, UnitPrice:50}}}
	promo := &entity.Promotion{Model: entity.Model{ID:1}, Type: "percentage", Discount: 10, StartDate: time.Now().Add(-time.Hour), EndDate: time.Now().Add(time.Hour), Published:true, UsageLimit: 5}

	// build repository.Repository with our mock components
	repoVal := repository.Repository{
		CartRepo:      &simpleCartRepo{cart: cart},
		PromotionRepo: &simplePromoRepo{promo: promo},
		OrderRepo:     &simpleOrderRepo{},
		AddressRepo:   nil,
	}
	logger, _ := zap.NewDevelopment()
	svc := &orderService{repo: repoVal, logger: logger}
	// ensure voucher code string exists
	code := "PROMO10"
	promo.VoucherCode = &code
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}

	res, err := svc.CreateOrder(context.Background(), req, 1)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if res.Total <= 0 { t.Fatalf("expected positive total, got %v", res.Total) }
}

func TestCreateOrder_InvalidVoucher(t *testing.T){
	cart := &entity.Cart{CustomerID:1, Items: []entity.CartItem{{ProductVariantID:1, Quantity:1, UnitPrice:100}}}
	repoVal := repository.Repository{CartRepo: &simpleCartRepo{cart: cart}, PromotionRepo: &simplePromoRepo{promo: nil}, OrderRepo: &simpleOrderRepo{}, AddressRepo: nil}
	logger, _ := zap.NewDevelopment()
	svc := &orderService{repo: repoVal, logger: logger}
	code := "INVALID"
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
	_, err := svc.CreateOrder(context.Background(), req, 1)
	if err == nil { t.Fatalf("expected error for invalid voucher") }
}

