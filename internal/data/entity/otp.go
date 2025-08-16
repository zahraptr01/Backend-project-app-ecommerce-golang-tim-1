package entity

import "time"

type AuthOTP struct {
	Model
	UserID    uint      `json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OtpCode   string    `json:"otp_code"`
	ExpiresAt time.Time `json:"expires_at"`
}
