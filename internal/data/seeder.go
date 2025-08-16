package data

import (
	"fmt"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()
		for i := range seeds {
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seeds[i]).Error
			if nil != err {
				name := reflect.TypeOf(seeds[i]).String()
				errorMessage := err.Error()
				return fmt.Errorf("%s seeder fail with %s", name, errorMessage)
			}
		}
		return nil
	})
}

// DataSeeds data
func dataSeeds() []interface{} {
	return []interface{}{
		entity.SeedUsers(),
		entity.SeedCustomers(),
		entity.SeedAddresses(),
		entity.SeedCarts(),
		entity.SeedOrders(),
	}
}
