package repository

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	List(ctx context.Context, page, limit int, search string) ([]entity.Category, int64, error)
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	Create(ctx context.Context, c *entity.Category) error
	Update(ctx context.Context, c *entity.Category) error
	Delete(ctx context.Context, id uint) error
	TogglePublished(ctx context.Context, id uint, published bool) error
	IsNameExists(ctx context.Context, name string, excludeID uint) (bool, error)
	CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error)
}

type categoryRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewCategoryRepository(DB *gorm.DB, log *zap.Logger) CategoryRepository {
	return &categoryRepositoryImpl{
		DB:  DB,
		Log: log,
	}
}

func (r *categoryRepositoryImpl) List(ctx context.Context, page, limit int, search string) ([]entity.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var cats []entity.Category
	q := r.DB.WithContext(ctx).Model(&entity.Category{})
	if search != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Limit(limit).Offset((page - 1) * limit).Find(&cats).Error; err != nil {
		return nil, 0, err
	}
	return cats, total, nil
}

func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	var c entity.Category
	if err := r.DB.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, c *entity.Category) error {
	return r.DB.WithContext(ctx).Create(c).Error
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, c *entity.Category) error {
	return r.DB.WithContext(ctx).Model(&entity.Category{}).
		Where("id = ?", c.ID).
		Updates(map[string]any{
			"name":      c.Name,
			"icon":      c.Icon,
			"published": c.Published,
		}).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// guard: cek masih dipakai product
	cnt, err := r.CountProductsByCategory(ctx, id)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("category is in use by products")
	}
	return r.DB.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

func (r *categoryRepositoryImpl) TogglePublished(ctx context.Context, id uint, published bool) error {
	return r.DB.WithContext(ctx).Model(&entity.Category{}).
		Where("id = ?", id).
		Update("published", published).Error
}

func (r *categoryRepositoryImpl) IsNameExists(ctx context.Context, name string, excludeID uint) (bool, error) {
	q := r.DB.WithContext(ctx).Model(&entity.Category{}).Where("LOWER(name)=LOWER(?)", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepositoryImpl) CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}
