package utils

import (
	"errors"
	"regexp"
	"strings"
)

var phoneRe = regexp.MustCompile(`^\+?\d{8,15}$`)

func SplitEmailOrPhone(v string) (email string, phone string, err error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", "", errors.New("email_or_phone is required")
	}
	if strings.Contains(v, "@") {
		if !strings.Contains(v, ".") {
			return "", "", errors.New("invalid email format")
		}
		return strings.ToLower(v), "", nil
	}

	normalized := strings.ReplaceAll(v, " ", "")
	if !phoneRe.MatchString(normalized) {
		return "", "", errors.New("invalid phone format")
	}
	return "", normalized, nil
}

func Deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
