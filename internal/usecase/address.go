package usecase

import (
	"context"
	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
)

type AddressService interface {
	CreateAddress(ctx context.Context, req dto.CreateAddressRequest, customerID uint) (*dto.AddressResponse, error)
	UpdateAddress(ctx context.Context, id uint, req dto.CreateAddressRequest, customerID uint) (*dto.AddressResponse, error)
	DeleteAddress(ctx context.Context, id uint, customerID uint) error
	ListAddresses(ctx context.Context, customerID uint) ([]dto.AddressResponse, error)
	SetDefaultAddress(ctx context.Context, customerID uint, addressID uint) error
}

type addressService struct {
	repo   repository.Repository
	logger *zap.Logger
}

func NewAddressService(repo repository.Repository, logger *zap.Logger) AddressService {
	return &addressService{repo: repo, logger: logger}
}

func (s *addressService) CreateAddress(ctx context.Context, req dto.CreateAddressRequest, customerID uint) (*dto.AddressResponse, error) {
	a := &entity.Address{
		CustomerID: customerID,
		Fullname:   req.Fullname,
		Email:      req.Email,
		Address:    req.Address,
	}
	if err := s.repo.AddressRepo.CreateAddress(ctx, a); err != nil {
		return nil, err
	}
	return &dto.AddressResponse{ID: a.ID, Fullname: a.Fullname, Email: a.Email, Address: a.Address}, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, id uint, req dto.CreateAddressRequest, customerID uint) (*dto.AddressResponse, error) {
	a, err := s.repo.AddressRepo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.CustomerID != customerID {
		return nil, err
	}
	a.Fullname = req.Fullname
	a.Email = req.Email
	a.Address = req.Address
	if err := s.repo.AddressRepo.UpdateAddress(ctx, a); err != nil {
		return nil, err
	}
	return &dto.AddressResponse{ID: a.ID, Fullname: a.Fullname, Email: a.Email, Address: a.Address}, nil
}

func (s *addressService) DeleteAddress(ctx context.Context, id uint, customerID uint) error {
	a, err := s.repo.AddressRepo.GetAddressByID(ctx, id)
	if err != nil {
		return err
	}
	if a.CustomerID != customerID {
		return nil
	}
	return s.repo.AddressRepo.DeleteAddress(ctx, id)
}

func (s *addressService) ListAddresses(ctx context.Context, customerID uint) ([]dto.AddressResponse, error) {
	addrs, err := s.repo.AddressRepo.ListAddressesByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.AddressResponse, 0, len(addrs))
	for _, a := range addrs {
		res = append(res, dto.AddressResponse{ID: a.ID, Fullname: a.Fullname, Email: a.Email, Address: a.Address})
	}
	return res, nil
}

func (s *addressService) SetDefaultAddress(ctx context.Context, customerID uint, addressID uint) error {
	return s.repo.AddressRepo.SetDefaultAddress(ctx, customerID, addressID)
}
