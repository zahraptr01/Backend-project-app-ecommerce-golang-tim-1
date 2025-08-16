package entity

import "time"

type Banner struct {
	Model
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	EndDate     time.Time `json:"end_date"`
	TargetURL   string    `json:"target_url"`
	BannerType  string    `json:"banner_type"`
	Image       string    `json:"image"`
	Published   bool      `json:"published"`
}
