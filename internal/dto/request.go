package dto

type LoginRequest struct {
	EmailOrPhone string `json:"email_or_phone" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ValidateOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=4"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type RegisterRequest struct {
	Fullname     string `json:"fullname" binding:"required,min=2"`
	EmailOrPhone string `json:"email_or_phone" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
}
