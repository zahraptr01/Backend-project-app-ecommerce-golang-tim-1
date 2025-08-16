package repository

import (
	"context"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	RedisRepo     RedisRepository
	AuthRepo      AuthRepository
	CustomerRepo  CustomerRepository
	StockRepo     StockRepository
	CategoryRepo  CategoryRepository
	BannerRepo    BannerRepository
	OrderRepo     OrderRepository
	AddressRepo   AddressRepository
	CartRepo      CartRepository
	PromotionRepo PromotionRepository
	UserRepo      UserRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) Repository {
	return Repository{
		RedisRepo:     NewRedisRepository(db, log),
		AuthRepo:      NewAuthRepository(db, log),
		CustomerRepo:  NewCustomerRepository(db, log),
		StockRepo:     NewStockRepository(db, log),
		CategoryRepo:  NewCategoryRepository(db, log),
		BannerRepo:    NewBannerRepository(db, log),
		OrderRepo:     NewOrderRepository(db, log),
		AddressRepo:   NewAddressRepository(db, log),
		CartRepo:      NewCartRepository(db, log),
		PromotionRepo: NewPromotionRepository(db, log),
		UserRepo:      NewUserRepository(db, log),
	}
}

// Repository interfaces for order and address
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByID(ctx context.Context, id uint) (*entity.Order, error)
	ListOrdersByCustomer(ctx context.Context, customerID uint, limit, offset int) ([]entity.Order, int64, error)
}

type AddressRepository interface {
	CreateAddress(ctx context.Context, addr *entity.Address) error
	UpdateAddress(ctx context.Context, addr *entity.Address) error
	DeleteAddress(ctx context.Context, id uint) error
	GetAddressByID(ctx context.Context, id uint) (*entity.Address, error)
	ListAddressesByCustomer(ctx context.Context, customerID uint) ([]entity.Address, error)
	SetDefaultAddress(ctx context.Context, customerID uint, addressID uint) error
}

// Cart repository
type CartRepository interface {
	GetCartByCustomer(ctx context.Context, customerID uint) (*entity.Cart, error)
	ClearCart(ctx context.Context, customerID uint) error
}

// Promotion repository
type PromotionRepository interface {
	GetByVoucherCode(ctx context.Context, code string) (*entity.Promotion, error)
	DecrementUsage(ctx context.Context, id uint) error
}
