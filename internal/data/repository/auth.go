package repository

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindCustomerByEmailOrPhone(ctx context.Context, identifier string) (*entity.User, error)
	SaveOTP(ctx context.Context, otp *entity.AuthOTP) error
	FindOTP(ctx context.Context, email string) (*entity.AuthOTP, error)
	DeleteOTP(ctx context.Context, email string) error
	UpdatePasswordByEmail(ctx context.Context, email string, newHashedPassword string) error
}

type authRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewAuthRepository(DB *gorm.DB, log *zap.Logger) AuthRepository {
	return &authRepositoryImpl{
		DB:  DB,
		Log: log,
	}
}

func (r *authRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) FindCustomerByEmailOrPhone(ctx context.Context, identifier string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).
		Where("email = ? OR phone = ?", identifier, identifier).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) SaveOTP(ctx context.Context, otp *entity.AuthOTP) error {
	return r.DB.WithContext(ctx).Create(otp).Error
}

func (r *authRepositoryImpl) FindOTP(ctx context.Context, email string) (*entity.AuthOTP, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	var otp entity.AuthOTP
	if err := r.DB.WithContext(ctx).Where("user_id = ?", user.ID).Order("created_at desc").First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepositoryImpl) DeleteOTP(ctx context.Context, email string) error {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Where("user_id = ?", user.ID).Delete(&entity.AuthOTP{}).Error
}

func (r *authRepositoryImpl) UpdatePasswordByEmail(ctx context.Context, email string, newHashedPassword string) error {
	result := r.DB.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Update("password", newHashedPassword)
	if result.RowsAffected == 0 {
		return errors.New("no user found with that email")
	}
	return result.Error
}
