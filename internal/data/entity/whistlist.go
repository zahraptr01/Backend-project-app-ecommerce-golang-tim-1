package entity

type Wishlist struct {
	Model
	CustomerID       uint           `json:"customer_id"`
	Customer         *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
}
