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

type CustomerService interface {
	RegisterCustomer(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type customerService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewCustomerService(repo repository.Repository, logger *zap.Logger, config utils.Configuration) CustomerService {
	return &customerService{
		Repo:   repo,
		Logger: logger,
		Config: config,
	}
}

func (s *customerService) RegisterCustomer(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {

	email, phone, err := utils.SplitEmailOrPhone(req.EmailOrPhone)
	if err != nil {
		return nil, err
	}

	if email != "" {
		exists, err := s.Repo.CustomerRepo.IsEmailExists(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already registered")
		}
	}
	if phone != "" {
		exists, err := s.Repo.CustomerRepo.IsPhoneExists(ctx, phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("phone already registered")
		}
	}

	hashed := utils.HashPassword(req.Password)

	var emailPtr, phonePtr *string
	if email != "" {
		emailPtr = &email
	}
	if phone != "" {
		phonePtr = &phone
	}

	user := &entity.User{
		Fullname: req.Fullname,
		Email:    emailPtr,
		Phone:    phonePtr,
		Password: hashed,
		Role:     "customer",
	}

	createdUser, _, err := s.Repo.CustomerRepo.CreateUserAndCustomer(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Fullname: createdUser.Fullname,
		Email:    utils.Deref(createdUser.Email),
		Phone:    utils.Deref(createdUser.Phone),
	}, nil
}
