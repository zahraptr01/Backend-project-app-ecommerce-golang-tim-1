package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/mailersend/mailersend-go"
)

func GenerateOTP(length int) string {
	digits := "0123456789"
	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			// fallback ke angka nol jika gagal
			otp[i] = '0'
		} else {
			otp[i] = digits[num.Int64()]
		}
	}
	return string(otp)
}

func SendOTPEmail(config Configuration, toEmail, otp string) error {
	apiKey := config.MailersendApiKey
	fromEmail := config.MailersendFromEmail
	fromName := config.AppName

	if apiKey == "" || fromEmail == "" || fromName == "" {
		return fmt.Errorf("mailersend env missing: MAILERSEND_API_KEY / MAILERSEND_FROM_EMAIL / MAILERSEND_FROM_NAME")
	}

	ms := mailersend.NewMailersend(apiKey)

	// Build message
	msg := ms.Email.NewMessage()
	msg.SetFrom(mailersend.From{Email: fromEmail, Name: fromName})
	msg.SetRecipients([]mailersend.Recipient{
		{Email: toEmail}, // Name optional
	})
	msg.SetSubject("Kode OTP Reset Password")
	msg.SetText(fmt.Sprintf("Kode OTP Anda: %s (berlaku 5 menit)", otp))
	msg.SetHTML(fmt.Sprintf(`
<!doctype html>
<html>
  <body style="font-family:Arial,sans-serif">
    <h2>Reset Password</h2>
    <p>Kode OTP Anda:</p>
    <div style="font-size:24px;font-weight:bold;letter-spacing:4px">%s</div>
    <p style="color:#666">Kode berlaku selama 5 menit.</p>
  </body>
</html>`, otp))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ms.Email.Send(ctx, msg)
	return err
}
