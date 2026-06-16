package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRToken генерирует уникальный токен для QR-кода
func GenerateQRToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback на timestamp если rand не сработал
		return fmt.Sprintf("token_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// GenerateQRCode генерирует QR-код в формате PNG
func GenerateQRCode(data string, size int) ([]byte, error) {
	qr, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации QR-кода: %w", err)
	}
	return qr, nil
}

// FormatDateTime форматирует дату и время в читаемый формат
func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

// FormatDateOnly форматирует только дату
func FormatDateOnly(t time.Time) string {
	return t.Format("02.01.2006")
}

// FormatTimeOnly форматирует только время
func FormatTimeOnly(t time.Time) string {
	return t.Format("15:04")
}

// IsRegistrationOpen проверяет, открыта ли регистрация
func IsRegistrationOpen(status string, maxParticipants *int, currentCount int) bool {
	if status != "open" {
		return false
	}
	if maxParticipants != nil && currentCount >= *maxParticipants {
		return false
	}
	return true
}

// ValidateEmail простая валидация email
func ValidateEmail(email string) bool {
	if len(email) < 5 {
		return false
	}
	for i, c := range email {
		if c == '@' && i > 0 && i < len(email)-1 {
			return true
		}
	}
	return false
}

// TruncateString обрезает строку до максимальной длины
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}