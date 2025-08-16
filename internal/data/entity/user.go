package entity

import "project-app-ecommerce-golang-tim-1/pkg/utils"

type User struct {
	Model
	Fullname  string     `json:"fullname"`
	Email     *string    `gorm:"uniqueIndex" json:"email,omitempty"`
	Phone     *string    `gorm:"uniqueIndex" json:"phone,omitempty"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	Customers []Customer `gorm:"foreignKey:UserID" json:"customers,omitempty"`
	AuthOTPs  []AuthOTP  `gorm:"foreignKey:UserID" json:"auth_otps,omitempty"`
}

type Customer struct {
	Model
	UserID    uint       `json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Addresses []Address  `gorm:"foreignKey:CustomerID" json:"addresses,omitempty"`
	Cart      Cart       `gorm:"foreignKey:CustomerID" json:"cart,omitempty"`
	Orders    []Order    `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
	Ratings   []Rating   `gorm:"foreignKey:CustomerID" json:"ratings,omitempty"`
	Wishlists []Wishlist `gorm:"foreignKey:CustomerID" json:"wishlists,omitempty"`
}

func SeedUsers() []User {
	email := "budi@example.com"
	users := []User{
		{
			Fullname: "Budi Santoso",
			Email:    &email,
			Password: utils.HashPassword("password123"),
			Role:     "superadmin",
		},
	}

	return users
}
