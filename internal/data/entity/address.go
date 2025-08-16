package entity

type Address struct {
	Model
	CustomerID uint      `json:"customer_id"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
}

// SeedAddresses returns default addresses for seeding
// func SeedAddresses() []Address {
// 	return []Address{
// 		{
// 			Fullname: "Budi Santoso",
// 			Email:    "budi@example.com",
// 			Address:  "Jl. Merdeka No.1, Jakarta",
// 		},
// 	}
// }

func SeedAddresses() []Address {
	return []Address{
		{
			CustomerID: 4,
			Fullname:   "Zahra",
			Email:      "zahra@example.com",
			Address:    "Jl. Kebangsaan No.10, Bandung",
			IsDefault:  true,
		},
	}
}
