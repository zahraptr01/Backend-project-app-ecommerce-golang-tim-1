//go:build integration
// +build integration

package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
)

func TestExpiredVoucherIntegration(t *testing.T){
	tdb := SetupTestDB(t)
	defer tdb.TearDown(t)

	// create expired promo
	now := time.Now()
	expired := entity.Promotion{Model: entity.Model{ID:10}, Name: "EXPIRED", Type: "fixed", Discount: 100, StartDate: now.Add(-48*time.Hour), EndDate: now.Add(-24*time.Hour), UsageLimit: 5, VoucherCode: func() *string { s := "EXPIRED"; return &s }(), Published: true}
	if err := tdb.DB.Create(&expired).Error; err != nil { t.Fatalf("failed to create expired promo: %v", err) }

	repo := repository.NewRepository(tdb.DB, zap.NewNop())
	svc := usecase.NewOrderService(repo, zap.NewNop())
	code := "EXPIRED"
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
	_, err := svc.CreateOrder(context.Background(), req, 1)
	if err == nil { t.Fatalf("expected error for expired voucher") }
}

func TestFixedVoucherIntegration(t *testing.T){
	tdb := SetupTestDB(t)
	defer tdb.TearDown(t)

	fixed := entity.Promotion{Model: entity.Model{ID:11}, Name: "FIXED", Type: "fixed", Discount: 50, StartDate: time.Now().Add(-time.Hour), EndDate: time.Now().Add(time.Hour), UsageLimit: 5, VoucherCode: func() *string { s := "FIXED"; return &s }(), Published: true}
	if err := tdb.DB.Create(&fixed).Error; err != nil { t.Fatalf("failed to create fixed promo: %v", err) }

	repo := repository.NewRepository(tdb.DB, zap.NewNop())
	svc := usecase.NewOrderService(repo, zap.NewNop())
	code := "FIXED"
	req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
	res, err := svc.CreateOrder(context.Background(), req, 1)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if res.Total <= 0 { t.Fatalf("expected positive total, got %v", res.Total) }
}

func TestConcurrentUsageIntegration(t *testing.T){
	tdb := SetupTestDB(t)
	defer tdb.TearDown(t)

	promo := entity.Promotion{Model: entity.Model{ID:20}, Name: "CONC", Type: "fixed", Discount: 1, StartDate: time.Now().Add(-time.Hour), EndDate: time.Now().Add(time.Hour), UsageLimit: 1, VoucherCode: func() *string { s := "CONC"; return &s }(), Published: true}
	if err := tdb.DB.Create(&promo).Error; err != nil { t.Fatalf("failed to create conc promo: %v", err) }

	repo := repository.NewRepository(tdb.DB, zap.NewNop())
	svc := usecase.NewOrderService(repo, zap.NewNop())
	code := "CONC"
	wg := sync.WaitGroup{}
	wg.Add(2)
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		go func(idx int){
			defer wg.Done()
			req := dto.CreateOrderRequest{AddressID: 1, PaymentMethod: "gopay", VoucherCode: &code}
			_, errs[idx] = svc.CreateOrder(context.Background(), req, 1)
		}(i)
	}
	wg.Wait()
	// one should succeed, one should fail due to usage limit
	if (errs[0] == nil && errs[1] == nil) || (errs[0] != nil && errs[1] != nil) {
		t.Fatalf("expected exactly one success and one failure, got errs: %v", errs)
	}
}
