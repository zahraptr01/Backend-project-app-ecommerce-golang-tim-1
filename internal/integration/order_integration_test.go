//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
)

func TestCreateOrderWithVoucherIntegration(t *testing.T){
	tdb := SetupTestDB(t)
	defer tdb.TearDown(t)

	// build repository using gorm DB
	repo := repository.NewRepository(tdb.DB, zap.NewNop())
	logger, _ := zap.NewDevelopment()
	svc := usecase.NewOrderService(repo, logger)

	// use seeded customer id 1 and voucher PROMO10
	code := "PROMO10"
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}

	_, err := svc.CreateOrder(context.Background(), req, 1)
	if err != nil {
		t.Fatalf("integration create order failed: %v", err)
	}
}
