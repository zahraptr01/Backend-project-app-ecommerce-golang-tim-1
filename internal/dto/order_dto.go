package dto

type CreateOrderRequest struct {
	AddressID     uint    `json:"address_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Note          *string `json:"note"`
	VoucherCode   *string `json:"voucher_code"`
}

type OrderItemDTO struct {
	ProductVariantID uint    `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
}

type OrderResponse struct {
	ID     uint           `json:"id"`
	Items  []OrderItemDTO `json:"items"`
	Total  float64        `json:"total"`
	Status string         `json:"status"`
}
