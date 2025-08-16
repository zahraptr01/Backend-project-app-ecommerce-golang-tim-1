//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"project-app-ecommerce-golang-tim-1/internal/data"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
)

// TestDB wraps testcontainer and gorm DB
type TestDB struct{
	DB *gorm.DB
	container testcontainers.Container
}

func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_pass",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	ip, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_pass dbname=test_db sslmode=disable", ip, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("failed to connect db: %v", err)
	}

	// run migrations
	if err := data.AutoMigrate(db); err != nil {
		container.Terminate(ctx)
		t.Fatalf("migration failed: %v", err)
	}

	// seed minimal necessary data
	if err := seedMinimal(db); err != nil {
		container.Terminate(ctx)
		t.Fatalf("seed failed: %v", err)
	}

	return &TestDB{DB: db, container: container}
}

func (tdb *TestDB) TearDown(t *testing.T) {
	ctx := context.Background()
	if err := tdb.container.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate container: %v", err)
	}
}

func seedMinimal(db *gorm.DB) error {
	// create a user + customer
	email := "test@example.com"
	u := entity.User{Fullname: "Test User", Email: &email, Password: "password", Role: "customer"}
	if err := db.Create(&u).Error; err != nil { return err }
	c := entity.Customer{UserID: u.ID}
	if err := db.Create(&c).Error; err != nil { return err }

	// address (linked to customer)
	a := entity.Address{CustomerID: c.ID, Fullname: "Test", Email: &email, Address: "Test St"}
	if err := db.Create(&a).Error; err != nil { return err }

	// product + variant
	p := entity.Product{Model: entity.Model{ID:1}, Name: "P1", SKU: "P1", Price: 100}
	if err := db.Create(&p).Error; err != nil { return err }
	v := entity.ProductVariant{Model: entity.Model{ID:1}, ProductID: 1, Variant: "V1", Stock: 10}
	if err := db.Create(&v).Error; err != nil { return err }

	// cart and item
	cart := entity.Cart{Model: entity.Model{ID:1}, CustomerID:1}
	if err := db.Create(&cart).Error; err != nil { return err }
	cartItem := entity.CartItem{CartID: cart.ID, ProductVariantID: v.ID, Quantity: 1, UnitPrice: 100}
	if err := db.Create(&cartItem).Error; err != nil { return err }

	// promotion
	now := time.Now()
	promo := entity.Promotion{Model: entity.Model{ID:1}, Name: "PROMO10", Type: "percentage", Discount: 10, StartDate: now.Add(-time.Hour), EndDate: now.Add(time.Hour), UsageLimit: 5, VoucherCode: func() *string { s := "PROMO10"; return &s }(), Published: true, ShowOnCheckout: true}
	if err := db.Create(&promo).Error; err != nil { return err }

	return nil
}

func TestMain(m *testing.M) {
	// keep TestMain minimal; tests will call SetupTestDB when needed
	os.Exit(m.Run())
}
