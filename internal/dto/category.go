package dto

type CategoryListQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type CategoryRow struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Published bool   `json:"published"`
}

type CategoryListResponse struct {
	Items        []CategoryRow `json:"items"`
	CurrentPage  int           `json:"current_page"`
	Limit        int           `json:"limit"`
	TotalPages   int           `json:"total_pages"`
	TotalRecords int64         `json:"total_records"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
	Icon string `json:"icon" binding:"required"` // path/url gambar
}

type UpdateCategoryRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,min=2"`
	Icon string `json:"icon" binding:"required"`
	// published di-toggle lewat endpoint khusus
}

type TogglePublishRequest struct {
	ID        uint `json:"id" binding:"required"`
	Published bool `json:"published"`
}
