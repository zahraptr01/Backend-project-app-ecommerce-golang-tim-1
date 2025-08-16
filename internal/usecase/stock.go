package usecase

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/pkg/utils"

	"go.uber.org/zap"
)

type StockService interface {
	List(ctx context.Context, q dto.StockListQuery) (*dto.StockListResponse, error)
	Detail(ctx context.Context, variantID uint) (*entity.ProductVariant, error)
	Add(ctx context.Context, req dto.AddStockRequest) error
	Set(ctx context.Context, req dto.SetStockRequest) error
	Delete(ctx context.Context, req dto.DeleteStockRequest) error
	VariantsDropdown(ctx context.Context, q dto.VariantDropdownQuery) (*dto.StockListResponse, error)
}

type stockService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewStockService(repo repository.Repository, logger *zap.Logger, config utils.Configuration) StockService {
	return &stockService{
		Repo:   repo,
		Logger: logger,
		Config: config,
	}
}

func (s *stockService) List(ctx context.Context, q dto.StockListQuery) (*dto.StockListResponse, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	rows, total, err := s.Repo.StockRepo.ListStock(ctx, q.Page, q.PageSize, q.Search)
	if err != nil {
		return nil, err
	}
	items := make([]dto.StockRow, len(rows))
	for i, r := range rows {
		items[i] = dto.StockRow(r)
	}

	totalPages := 0
	if q.PageSize > 0 {
		// ceiling division
		totalPages = int((total + int64(q.PageSize) - 1) / int64(q.PageSize))
	}

	return &dto.StockListResponse{
		Items:        items,
		CurrentPage:  q.Page,
		Limit:        q.PageSize,
		TotalPages:   totalPages,
		TotalRecords: total,
	}, nil
}

func (s *stockService) Detail(ctx context.Context, variantID uint) (*entity.ProductVariant, error) {
	if variantID == 0 {
		return nil, errors.New("variant_id required")
	}
	return s.Repo.StockRepo.GetVariantStock(ctx, variantID)
}

func (s *stockService) Add(ctx context.Context, req dto.AddStockRequest) error {
	return s.Repo.StockRepo.IncreaseStock(ctx, req.VariantID, req.AddQty)
}

func (s *stockService) Set(ctx context.Context, req dto.SetStockRequest) error {
	return s.Repo.StockRepo.SetStock(ctx, req.VariantID, req.Qty)
}

func (s *stockService) Delete(ctx context.Context, req dto.DeleteStockRequest) error {
	return s.Repo.StockRepo.DeleteStock(ctx, req.VariantID)
}

func (s *stockService) VariantsDropdown(ctx context.Context, q dto.VariantDropdownQuery) (*dto.StockListResponse, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	rows, total, err := s.Repo.StockRepo.ListVariantsForDropdown(ctx, q.Page, q.PageSize, q.Search)
	if err != nil {
		return nil, err
	}
	items := make([]dto.StockRow, len(rows))
	for i, r := range rows {
		items[i] = dto.StockRow(r)
	}

	totalPages := 0
	if q.PageSize > 0 {
		totalPages = int((total + int64(q.PageSize) - 1) / int64(q.PageSize))
	}

	return &dto.StockListResponse{
		Items:        items,
		CurrentPage:  q.Page,
		Limit:        q.PageSize,
		TotalPages:   totalPages,
		TotalRecords: total,
	}, nil
}
