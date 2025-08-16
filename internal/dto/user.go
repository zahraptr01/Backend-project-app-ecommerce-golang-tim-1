package dto

type CreateUserRequest struct {
	Fullname string  `json:"fullname" binding:"required"`
	Email    *string `json:"email" binding:"required,email"`
	Phone    *string `json:"phone" binding:"required"`
	Role     string  `json:"role" binding:"required"` // admin atau staff
}

type UpdateUserRequest struct {
	Fullname string  `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Role     string  `json:"role,omitempty"` // admin atau staff
	IsActive *bool   `json:"is_active,omitempty"`
}

type UserResponse struct {
	ID        uint    `json:"id"`
	Fullname  string  `json:"fullname"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
