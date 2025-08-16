package data

import (
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Customer{},
		&entity.Address{},
		&entity.Category{},
		&entity.Product{},
		&entity.ProductVariant{},
		&entity.ProductPhoto{},
		&entity.Cart{},
		&entity.CartItem{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Wishlist{},
		&entity.Rating{},
		&entity.Promotion{},
		&entity.PromotionProduct{},
		&entity.Banner{},
		&entity.AuthOTP{},
	)
}
