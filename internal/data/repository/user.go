package repository

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindAll(ctx context.Context, sort, order string, page, limit int) ([]entity.User, error)
	FindByID(ctx context.Context, id uint) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id uint) error
}

type userRepositoryImpl struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) UserRepository {
	return &userRepositoryImpl{db: db, log: log}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, sort, order string, page, limit int) ([]entity.User, error) {
	var users []entity.User
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("role IN ?", []string{"admin", "superadmin"}).
		Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ? AND role IN ?", id, []string{"admin", "superadmin"}).Error
	return user, err
}

func (r *userRepositoryImpl) Update(ctx context.Context, user entity.User) (entity.User, error) {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if user.Fullname != "" {
		updates["fullname"] = user.Fullname
	}
	if user.Email != nil {
		updates["email"] = *user.Email
	}
	if user.Phone != nil {
		updates["phone"] = *user.Phone
	}
	if user.Role != "" {
		updates["role"] = user.Role
	}
	updates["is_active"] = user.IsActive

	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND role IN ?", user.ID, []string{"admin", "staff"}).
		Updates(updates).Error

	return user, err
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ? AND role IN ?", id, []string{"admin", "superadmin"}).Error
}
