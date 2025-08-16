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

type CategoryService interface {
	List(ctx context.Context, q dto.CategoryListQuery) (*dto.CategoryListResponse, error)
	Get(ctx context.Context, id uint) (*entity.Category, error)
	Create(ctx context.Context, req dto.CreateCategoryRequest) error
	Update(ctx context.Context, req dto.UpdateCategoryRequest) error
	Delete(ctx context.Context, id uint) error
	TogglePublished(ctx context.Context, req dto.TogglePublishRequest) error
}

type categoryService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewCategoryService(repo repository.Repository, logger *zap.Logger, config utils.Configuration) CategoryService {
	return &categoryService{
		Repo:   repo,
		Logger: logger,
		Config: config,
	}
}

func (s *categoryService) List(ctx context.Context, q dto.CategoryListQuery) (*dto.CategoryListResponse, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}

	rows, total, err := s.Repo.CategoryRepo.List(ctx, q.Page, q.Limit, q.Search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.CategoryRow, len(rows))
	for i, c := range rows {
		items[i] = dto.CategoryRow{
			ID: c.ID, Name: c.Name, Icon: c.Icon, Published: c.Published,
		}
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	return &dto.CategoryListResponse{
		Items: items, CurrentPage: q.Page, Limit: q.Limit,
		TotalPages: totalPages, TotalRecords: total,
	}, nil
}

func (s *categoryService) Get(ctx context.Context, id uint) (*entity.Category, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.Repo.CategoryRepo.GetByID(ctx, id)
}

func (s *categoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) error {
	exists, err := s.Repo.CategoryRepo.IsNameExists(ctx, req.Name, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category name already exists")
	}

	return s.Repo.CategoryRepo.Create(ctx, &entity.Category{
		Name: req.Name, Icon: req.Icon, Published: true, // default published true
	})
}

func (s *categoryService) Update(ctx context.Context, req dto.UpdateCategoryRequest) error {
	if req.ID == 0 {
		return errors.New("invalid id")
	}
	exists, err := s.Repo.CategoryRepo.IsNameExists(ctx, req.Name, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category name already exists")
	}

	return s.Repo.CategoryRepo.Update(ctx, &entity.Category{
		Model: entity.Model{ID: req.ID}, Name: req.Name, Icon: req.Icon,
	})
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return s.Repo.CategoryRepo.Delete(ctx, id)
}

func (s *categoryService) TogglePublished(ctx context.Context, req dto.TogglePublishRequest) error {
	if req.ID == 0 {
		return errors.New("invalid id")
	}
	return s.Repo.CategoryRepo.TogglePublished(ctx, req.ID, req.Published)
}
