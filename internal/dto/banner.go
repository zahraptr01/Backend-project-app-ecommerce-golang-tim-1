package dto

import "time"

type BannerListQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
}

type BannerRow struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	EndDate     time.Time `json:"end_date"`
	TargetURL   string    `json:"target_url"`
	BannerType  string    `json:"banner_type"`
	Image       string    `json:"image"`
	Published   bool      `json:"published"`
}

type BannerListResponse struct {
	Items        []BannerRow `json:"items"`
	CurrentPage  int         `json:"current_page"`
	Limit        int         `json:"limit"`
	TotalPages   int         `json:"total_pages"`
	TotalRecords int64       `json:"total_records"`
}

type CreateBannerRequest struct {
	Name        string    `form:"name" binding:"required,min=2"`
	ReleaseDate time.Time `form:"release_date" time_format:"2006-01-02"`
	EndDate     time.Time `form:"end_date" time_format:"2006-01-02"`
	TargetURL   string    `form:"target_url" binding:"required"`
	BannerType  string    `form:"banner_type" binding:"required"`
	Image       string    `form:"image"`
	Published   bool      `form:"published"`
}

type UpdateBannerRequest struct {
	ID          uint      `form:"id" binding:"required"`
	Name        string    `form:"name" binding:"required,min=2"`
	ReleaseDate time.Time `form:"release_date" time_format:"2006-01-02"`
	EndDate     time.Time `form:"end_date" time_format:"2006-01-02"`
	TargetURL   string    `form:"target_url" binding:"required"`
	BannerType  string    `form:"banner_type" binding:"required"`
	Image       string    `form:"image"`
	Published   *bool     `form:"published"`
}

type ToggleBannerPublishRequest struct {
	ID        uint `json:"id" binding:"required"`
	Published bool `json:"published"`
}
