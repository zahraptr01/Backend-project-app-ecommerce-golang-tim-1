package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

type promotionRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewPromotionRepository(db *gorm.DB, log *zap.Logger) PromotionRepository {
	return &promotionRepo{db: db, log: log}
}

func (r *promotionRepo) GetByVoucherCode(ctx context.Context, code string) (*entity.Promotion, error) {
	var p entity.Promotion
	if err := r.db.WithContext(ctx).Where("voucher_code = ?", code).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *promotionRepo) DecrementUsage(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Model(&entity.Promotion{}).Where("id = ? AND usage_limit > 0", id).UpdateColumn("usage_limit", gorm.Expr("usage_limit - 1"))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("usage limit exceeded")
	}
	// touch updated_at
	r.db.WithContext(ctx).Model(&entity.Promotion{}).Where("id = ?", id).UpdateColumn("updated_at", clause.Expr{SQL: "NOW()"})
	return nil
}
