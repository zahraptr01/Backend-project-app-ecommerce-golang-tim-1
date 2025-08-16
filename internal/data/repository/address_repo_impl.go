package repository

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

type addressRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewAddressRepository(db *gorm.DB, log *zap.Logger) AddressRepository {
	return &addressRepo{db: db, log: log}
}

func (r *addressRepo) CreateAddress(ctx context.Context, addr *entity.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *addressRepo) UpdateAddress(ctx context.Context, addr *entity.Address) error {
	return r.db.WithContext(ctx).Save(addr).Error
}

func (r *addressRepo) DeleteAddress(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Address{}, id).Error
}

func (r *addressRepo) GetAddressByID(ctx context.Context, id uint) (*entity.Address, error) {
	var a entity.Address
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *addressRepo) ListAddressesByCustomer(ctx context.Context, customerID uint) ([]entity.Address, error) {
	var addrs []entity.Address
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&addrs).Error; err != nil {
		return nil, err
	}
	return addrs, nil
}

func (r *addressRepo) SetDefaultAddress(ctx context.Context, customerID uint, addressID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ensure the address belongs to the customer
	var addr entity.Address
	if err := tx.Where("id = ? AND customer_id = ?", addressID, customerID).First(&addr).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Unset previous defaults
	if err := tx.Model(&entity.Address{}).Where("customer_id = ? AND is_default = ?", customerID, true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the selected address as default
	if err := tx.Model(&entity.Address{}).Where("id = ?", addressID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
