package repository

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockRow struct {
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	VariantID    uint   `json:"variant_id"`
	VariantName  string `json:"variant_name"`
	Quantity     int    `json:"quantity"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type StockRepository interface {
	// List stok dengan pagination + optional pencarian by product/variant
	ListStock(ctx context.Context, page, pageSize int, search string) ([]StockRow, int64, error)

	// Detail stok per variant
	GetVariantStock(ctx context.Context, variantID uint) (*entity.ProductVariant, error)

	// Tambah stok (delta +n)
	IncreaseStock(ctx context.Context, variantID uint, addQty int) error

	// Edit stok (set absolut = qty)
	SetStock(ctx context.Context, variantID uint, qty int) error

	// Delete stok (set ke 0)
	DeleteStock(ctx context.Context, variantID uint) error

	// Util untuk dropdown: ambil variant beserta nama produk (pagination + search optional)
	ListVariantsForDropdown(ctx context.Context, page, pageSize int, search string) ([]StockRow, int64, error)

	// Untuk order: kurangi stok (validate tidak boleh minus)
	DecreaseStock(ctx context.Context, variantID uint, qty int) error
}

type stockRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewStockRepository(DB *gorm.DB, log *zap.Logger) StockRepository {
	return &stockRepositoryImpl{
		DB:  DB,
		Log: log,
	}
}

func (r *stockRepositoryImpl) ListStock(ctx context.Context, page, pageSize int, search string) ([]StockRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var rows []StockRow
	q := r.DB.WithContext(ctx).
		Table("product_variants pv").
		Select(`p.id as product_id, p.name as product_name,
		        pv.id as variant_id, pv.variant as variant_name,
		        pv.stock as quantity,
		        p.category_id as category_id, c.name as category_name`).
		Joins("JOIN products p ON p.id = pv.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id")

	if search != "" {
		q = q.Where("LOWER(p.name) LIKE LOWER(?) OR LOWER(pv.variant) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Order("p.name ASC, pv.variant ASC").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *stockRepositoryImpl) GetVariantStock(ctx context.Context, variantID uint) (*entity.ProductVariant, error) {
	var v entity.ProductVariant
	if err := r.DB.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		First(&v, "id = ?", variantID).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *stockRepositoryImpl) IncreaseStock(ctx context.Context, variantID uint, addQty int) error {
	if addQty <= 0 {
		return errors.New("addQty must be > 0")
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var v entity.ProductVariant
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}). // <-- fix
			First(&v, "id = ?", variantID).Error; err != nil {
			return err
		}
		v.Stock += addQty
		return tx.Model(&entity.ProductVariant{}).
			Where("id = ?", variantID).
			Update("stock", v.Stock).Error
	})
}

func (r *stockRepositoryImpl) SetStock(ctx context.Context, variantID uint, qty int) error {
	if qty < 0 {
		return errors.New("qty must be >= 0")
	}
	return r.DB.WithContext(ctx).Model(&entity.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", qty).Error
}

func (r *stockRepositoryImpl) DeleteStock(ctx context.Context, variantID uint) error {
	// definisi "delete" = set 0
	return r.DB.WithContext(ctx).Model(&entity.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", 0).Error
}

func (r *stockRepositoryImpl) ListVariantsForDropdown(ctx context.Context, page, pageSize int, search string) ([]StockRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var rows []StockRow
	q := r.DB.WithContext(ctx).
		Table("product_variants pv").
		Select(`p.id as product_id, p.name as product_name,
		        pv.id as variant_id, pv.variant as variant_name,
		        pv.stock as quantity,
		        p.category_id as category_id, c.name as category_name`).
		Joins("JOIN products p ON p.id = pv.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id")

	if search != "" {
		q = q.Where("LOWER(p.name) LIKE LOWER(?) OR LOWER(pv.variant) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Order("p.name ASC, pv.variant ASC").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *stockRepositoryImpl) DecreaseStock(ctx context.Context, variantID uint, qty int) error {
	if qty <= 0 {
		return errors.New("qty must be > 0")
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var v entity.ProductVariant
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}). // <-- fix
			First(&v, "id = ?", variantID).Error; err != nil {
			return err
		}
		if v.Stock < qty {
			return errors.New("insufficient stock")
		}
		newQty := v.Stock - qty
		return tx.Model(&entity.ProductVariant{}).
			Where("id = ?", variantID).
			Update("stock", newQty).Error
	})
}
