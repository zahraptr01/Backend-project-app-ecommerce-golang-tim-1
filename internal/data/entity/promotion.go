package entity

import "time"

type Promotion struct {
	Model
	Name           string             `json:"name"`
	Type           string             `json:"type"`
	StartDate      time.Time          `json:"start_date"`
	EndDate        time.Time          `json:"end_date"`
	Discount       float64            `json:"discount"`
	UsageLimit     int                `json:"usage_limit"`
	VoucherCode    *string            `json:"voucher_code,omitempty"`
	ShowOnCheckout bool               `json:"show_on_checkout"`
	Published      bool               `json:"published"`
	Products       []PromotionProduct `gorm:"foreignKey:PromotionID" json:"products,omitempty"`
}

type PromotionProduct struct {
	Model
	PromotionID uint      `json:"promotion_id"`
	ProductID   uint      `json:"product_id"`
	Promotion   Promotion `gorm:"foreignKey:PromotionID" json:"promotion,omitempty"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
