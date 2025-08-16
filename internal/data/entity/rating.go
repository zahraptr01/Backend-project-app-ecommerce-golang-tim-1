package entity

type Rating struct {
	Model
	CustomerID uint      `json:"customer_id"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Rating     int       `json:"rating"`
	Review     string    `json:"review"`
}
