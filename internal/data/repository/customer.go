package repository

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	CreateUserAndCustomer(ctx context.Context, user *entity.User,
	) (*entity.User, *entity.Customer, error)
}

type customerRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewCustomerRepository(DB *gorm.DB, log *zap.Logger) CustomerRepository {
	return &customerRepositoryImpl{
		DB:  DB,
		Log: log,
	}
}

func (r *customerRepositoryImpl) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, nil
	}
	var count int64
	if err := r.DB.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *customerRepositoryImpl) IsPhoneExists(ctx context.Context, phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}
	var count int64
	if err := r.DB.WithContext(ctx).Model(&entity.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *customerRepositoryImpl) CreateUserAndCustomer(ctx context.Context, user *entity.User,
) (*entity.User, *entity.Customer, error) {

	if user == nil {
		return nil, nil, errors.New("user is nil")
	}

	returnedUser := new(entity.User)
	returnedCustomer := new(entity.Customer)

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create user
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		*returnedUser = *user

		// create customer
		customer := &entity.Customer{
			UserID: user.ID,
		}
		if err := tx.Create(customer).Error; err != nil {
			return err
		}
		*returnedCustomer = *customer
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return returnedUser, returnedCustomer, nil
}
