package dto

type CartItemResponse struct {
	ProductVariantID uint    `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
}

type CartResponse struct {
	CustomerID uint               `json:"customer_id"`
	Items      []CartItemResponse `json:"items"`
	Total      float64            `json:"total"`
}
