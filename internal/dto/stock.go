package dto

type StockListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Search   string `form:"search"`
}

type StockRow struct {
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	VariantID    uint   `json:"variant_id"`
	VariantName  string `json:"variant_name"`
	Quantity     int    `json:"quantity"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type StockListResponse struct {
	Items        []StockRow `json:"items"`
	CurrentPage  int        `json:"current_page"`
	Limit        int        `json:"limit"`
	TotalPages   int        `json:"total_pages"`
	TotalRecords int64      `json:"total_records"`
}

type AddStockRequest struct {
	VariantID uint `json:"variant_id" binding:"required"`
	AddQty    int  `json:"add_qty" binding:"required,gt=0"`
}

type SetStockRequest struct {
	VariantID uint `json:"variant_id" binding:"required"`
	Qty       int  `json:"qty" binding:"required,gte=0"`
}

type DeleteStockRequest struct {
	VariantID uint `json:"variant_id" binding:"required"`
}

type VariantDropdownQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Search   string `form:"search"`
}
