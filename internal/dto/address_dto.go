package dto

type CreateAddressRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Address  string `json:"address" binding:"required"`
}

type AddressResponse struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}
