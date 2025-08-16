package repository

import (
	"context"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BannerRepository interface {
	List(ctx context.Context, page, limit int, search, status string) ([]entity.Banner, int64, error)
	GetByID(ctx context.Context, id uint) (*entity.Banner, error)
	Create(ctx context.Context, b *entity.Banner) error
	Update(ctx context.Context, b *entity.Banner) error
	Delete(ctx context.Context, id uint) error
	TogglePublished(ctx context.Context, id uint, published bool) error
	AutoUnpublishExpired(ctx context.Context) (int64, error)
}

type bannerRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewBannerRepository(DB *gorm.DB, log *zap.Logger) BannerRepository {
	return &bannerRepositoryImpl{
		DB:  DB,
		Log: log,
	}
}

func (r *bannerRepositoryImpl) List(ctx context.Context, page, limit int, search, status string) ([]entity.Banner, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	var rows []entity.Banner
	var total int64

	q := r.DB.WithContext(ctx).Model(&entity.Banner{})
	if search != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?) OR LOWER(banner_type) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}
	now := time.Now()
	switch status {
	case "active":
		q = q.Where("release_date <= ? AND end_date >= ?", now, now)
	case "expired":
		q = q.Where("end_date < ?", now)
	case "upcoming":
		q = q.Where("release_date > ?", now)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Order("id DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *bannerRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Banner, error) {
	var b entity.Banner
	if err := r.DB.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *bannerRepositoryImpl) Create(ctx context.Context, b *entity.Banner) error {
	return r.DB.WithContext(ctx).Create(b).Error
}

func (r *bannerRepositoryImpl) Update(ctx context.Context, b *entity.Banner) error {
	return r.DB.WithContext(ctx).Model(&entity.Banner{}).
		Where("id = ?", b.ID).
		Updates(map[string]any{
			"name":         b.Name,
			"release_date": b.ReleaseDate,
			"end_date":     b.EndDate,
			"target_url":   b.TargetURL,
			"banner_type":  b.BannerType,
			"image":        b.Image,
			"published":    b.Published,
		}).Error
}

func (r *bannerRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&entity.Banner{}, id).Error
}

func (r *bannerRepositoryImpl) TogglePublished(ctx context.Context, id uint, published bool) error {
	return r.DB.WithContext(ctx).Model(&entity.Banner{}).
		Where("id = ?", id).
		Update("published", published).Error
}

func (r *bannerRepositoryImpl) AutoUnpublishExpired(ctx context.Context) (int64, error) {
	now := time.Now()
	res := r.DB.WithContext(ctx).Model(&entity.Banner{}).
		Where("end_date < ? AND published = ?", now, true).
		Update("published", false)
	return res.RowsAffected, res.Error
}
