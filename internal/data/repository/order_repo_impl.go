package repository

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

type orderRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewOrderRepository(db *gorm.DB, log *zap.Logger) OrderRepository {
	return &orderRepo{db: db, log: log}
}

func (r *orderRepo) CreateOrder(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepo) GetOrderByID(ctx context.Context, id uint) (*entity.Order, error) {
	var o entity.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepo) ListOrdersByCustomer(ctx context.Context, customerID uint, limit, offset int) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Order{}).Where("customer_id = ?", customerID)
	q.Count(&total)
	if err := q.Preload("Items").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
