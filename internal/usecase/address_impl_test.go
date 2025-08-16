package usecase

import (
    "context"
    "testing"

    "go.uber.org/zap"
    "project-app-ecommerce-golang-tim-1/internal/data/entity"
    "project-app-ecommerce-golang-tim-1/internal/data/repository"
)

type mockAddressRepo struct{}
func (r *mockAddressRepo) CreateAddress(ctx context.Context, addr *entity.Address) error { return nil }
func (r *mockAddressRepo) UpdateAddress(ctx context.Context, addr *entity.Address) error { return nil }
func (r *mockAddressRepo) DeleteAddress(ctx context.Context, id uint) error { return nil }
func (r *mockAddressRepo) GetAddressByID(ctx context.Context, id uint) (*entity.Address, error) { return &entity.Address{Model: entity.Model{ID:id}, CustomerID:1}, nil }
func (r *mockAddressRepo) ListAddressesByCustomer(ctx context.Context, customerID uint) ([]entity.Address, error) { return []entity.Address{{Model: entity.Model{ID:1}, CustomerID:customerID}}, nil }
func (r *mockAddressRepo) SetDefaultAddress(ctx context.Context, customerID uint, addressID uint) error { return nil }

func TestSetDefaultAddress(t *testing.T){
    repoVal := repository.Repository{AddressRepo: &mockAddressRepo{}}
    logger, _ := zap.NewDevelopment()
    svc := NewAddressService(repoVal, logger)
    if err := svc.SetDefaultAddress(context.Background(), 1, 1); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
