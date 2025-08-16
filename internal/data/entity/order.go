package entity

type Order struct {
	Model
	CustomerID     uint        `json:"customer_id"`
	Customer       *Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	AddressID      uint        `json:"address_id"`
	Address        Address     `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Note           string      `json:"note"`
	PaymentMethod  string      `json:"payment_method"`
	VoucherCode    *string     `json:"voucher_code"`
	PromotionID    *uint       `json:"promotion_id,omitempty"`
	Discount       float64     `json:"discount"`
	Status         string      `json:"status"`
	TrackingNumber *string     `json:"tracking_number"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	Model
	OrderID          uint           `json:"order_id"`
	Order            Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity         int            `json:"quantity"`
	UnitPrice        float64        `json:"unit_price"`
}

// SeedOrders returns default orders for seeding
func SeedOrders() []Order {
	return []Order{
		{
			CustomerID:    1,
			AddressID:     1,
			Note:          "Test order",
			PaymentMethod: "gopay",
			Status:        "created",
			Items: []OrderItem{
				{ProductVariantID: 1, Quantity: 1, UnitPrice: 100.0},
			},
		},
	}
}
