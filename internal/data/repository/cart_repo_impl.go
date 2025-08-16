package repository

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

type cartRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCartRepository(db *gorm.DB, log *zap.Logger) CartRepository {
	return &cartRepo{db: db, log: log}
}

func (r *cartRepo) GetCartByCustomer(ctx context.Context, customerID uint) (*entity.Cart, error) {
	var cart entity.Cart
	if err := r.db.WithContext(ctx).Preload("Items").Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) ClearCart(ctx context.Context, customerID uint) error {
	// delete cart items for customer
	if err := r.db.WithContext(ctx).Where("cart_id IN (SELECT id FROM carts WHERE customer_id = ?)", customerID).Delete(&entity.CartItem{}).Error; err != nil {
		return err
	}
	return nil
}
