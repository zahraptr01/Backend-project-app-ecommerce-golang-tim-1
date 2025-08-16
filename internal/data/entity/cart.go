package entity

type Cart struct {
	Model
	CustomerID uint       `json:"customer_id"`
	Customer   *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Items      []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
}

type CartItem struct {
	Model
	CartID           uint           `json:"cart_id"`
	Cart             Cart           `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity         int            `json:"quantity"`
	UnitPrice        float64        `json:"unit_price"`
}

// SeedCarts returns default carts for seeding
func SeedCarts() []Cart {
	return []Cart{
		{
			CustomerID: 1,
			Items: []CartItem{
				{ProductVariantID: 1, Quantity: 1, UnitPrice: 100.0},
			},
		},
	}
}
