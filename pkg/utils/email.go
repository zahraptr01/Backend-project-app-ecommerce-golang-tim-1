package utils

import (
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type smtpSender struct {
	host     string
	port     int
	email    string
	password string
}

func NewEmailSender(host string, port int, email, password string) EmailSender {
	return &smtpSender{host: host, port: port, email: email, password: password}
}

func (s *smtpSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.host, s.port, s.email, s.password)
	return d.DialAndSend(m)
}
