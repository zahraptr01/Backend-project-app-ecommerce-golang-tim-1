package entity

type Category struct {
	Model
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Published bool      `json:"published"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
