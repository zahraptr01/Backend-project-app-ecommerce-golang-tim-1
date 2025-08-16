package entity

type Product struct {
	Model
	Name        string             `json:"name"`
	SKU         string             `gorm:"unique" json:"sku"`
	CategoryID  uint               `json:"category_id"`
	Category    Category           `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Price       float64            `json:"price"`
	Description string             `json:"description"`
	Published   bool               `json:"published"`
	Variants    []ProductVariant   `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Photos      []ProductPhoto     `gorm:"foreignKey:ProductID" json:"photos,omitempty"`
	Promotions  []PromotionProduct `gorm:"foreignKey:ProductID" json:"promotions,omitempty"`
}

type ProductVariant struct {
	Model
	ProductID  uint        `json:"product_id"`
	Product    Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Variant    string      `json:"variant"`
	Stock      int         `json:"stock"`
	CartItems  []CartItem  `gorm:"foreignKey:ProductVariantID" json:"cart_items,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductVariantID" json:"order_items,omitempty"`
	Wishlists  []Wishlist  `gorm:"foreignKey:ProductVariantID" json:"wishlists,omitempty"`
}

type ProductPhoto struct {
	Model
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	URL       string  `json:"url"`
	IsDefault bool    `json:"is_default"`
}
