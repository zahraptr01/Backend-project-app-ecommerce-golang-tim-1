package usecase

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"time"

	"go.uber.org/zap"
)

type BannerService interface {
	List(ctx context.Context, q dto.BannerListQuery) (*dto.BannerListResponse, error)
	Get(ctx context.Context, id uint) (*entity.Banner, error)
	Create(ctx context.Context, req dto.CreateBannerRequest) error
	Update(ctx context.Context, req dto.UpdateBannerRequest) error
	Delete(ctx context.Context, id uint) error
	TogglePublished(ctx context.Context, id uint, published bool) error
}

type bannerService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewBannerService(repo repository.Repository, logger *zap.Logger, config utils.Configuration) BannerService {
	return &bannerService{
		Repo:   repo,
		Logger: logger,
		Config: config,
	}
}

func (s *bannerService) List(ctx context.Context, q dto.BannerListQuery) (*dto.BannerListResponse, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}

	_, _ = s.Repo.BannerRepo.AutoUnpublishExpired(ctx)

	rows, total, err := s.Repo.BannerRepo.List(ctx, q.Page, q.Limit, q.Search, q.Status)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BannerRow, len(rows))
	for i, b := range rows {
		items[i] = dto.BannerRow{
			ID: b.ID, Name: b.Name, ReleaseDate: b.ReleaseDate, EndDate: b.EndDate,
			TargetURL: b.TargetURL, BannerType: b.BannerType, Image: b.Image, Published: b.Published,
		}
	}
	totalPages := int((total + int64(q.Limit) - 1) / int64(q.Limit))

	return &dto.BannerListResponse{
		Items:        items,
		CurrentPage:  q.Page,
		Limit:        q.Limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}, nil
}

func (s *bannerService) Get(ctx context.Context, id uint) (*entity.Banner, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.Repo.BannerRepo.GetByID(ctx, id)
}

func (s *bannerService) Create(ctx context.Context, req dto.CreateBannerRequest) error {
	if !req.ReleaseDate.IsZero() && !req.EndDate.IsZero() && req.ReleaseDate.After(req.EndDate) {
		return errors.New("release_date must be before end_date")
	}
	b := &entity.Banner{
		Name:        req.Name,
		ReleaseDate: req.ReleaseDate,
		EndDate:     req.EndDate,
		TargetURL:   req.TargetURL,
		BannerType:  req.BannerType,
		Image:       req.Image,
		Published:   req.Published,
	}
	return s.Repo.BannerRepo.Create(ctx, b)
}

func (s *bannerService) Update(ctx context.Context, req dto.UpdateBannerRequest) error {
	if req.ID == 0 {
		return errors.New("invalid id")
	}
	ex, err := s.Repo.BannerRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	img := ex.Image
	if req.Image != "" {
		img = req.Image
	}
	published := ex.Published
	if req.Published != nil {
		published = *req.Published
	}
	if !req.ReleaseDate.IsZero() && !req.EndDate.IsZero() && req.ReleaseDate.After(req.EndDate) {
		return errors.New("release_date must be before end_date")
	}
	now := time.Now()
	if !req.EndDate.IsZero() && req.EndDate.Before(now) {
		published = false
	}

	return s.Repo.BannerRepo.Update(ctx, &entity.Banner{
		Model:       entity.Model{ID: req.ID},
		Name:        req.Name,
		ReleaseDate: req.ReleaseDate,
		EndDate:     req.EndDate,
		TargetURL:   req.TargetURL,
		BannerType:  req.BannerType,
		Image:       img,
		Published:   published,
	})
}

func (s *bannerService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return s.Repo.BannerRepo.Delete(ctx, id)
}

func (s *bannerService) TogglePublished(ctx context.Context, id uint, published bool) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	ex, err := s.Repo.BannerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !ex.EndDate.IsZero() && ex.EndDate.Before(time.Now()) {
		published = false
	}
	return s.Repo.BannerRepo.TogglePublished(ctx, id, published)
}
