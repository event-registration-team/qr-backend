package service

import (
	"bytes"
	"encoding/base64"
	"event-registration/internal/config"
	"event-registration/internal/models"
	"event-registration/pkg/utils"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	cfg config.SMTPConfig
}

func NewEmailService(cfg config.SMTPConfig) *EmailService {
	return &EmailService{cfg: cfg}
}

// SendRegistrationEmail отправляет письмо с QR-кодом после регистрации
func (s *EmailService) SendRegistrationEmail(participant *models.Participant, event *models.Event) error {
	// Генерируем QR-код
	qrPNG, err := utils.GenerateQRCode(participant.QRToken, 256)
	if err != nil {
		return fmt.Errorf("ошибка генерации QR-кода: %w", err)
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrPNG)

	subject := fmt.Sprintf("Вы зарегистрированы: %s", event.Title)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Вы успешно зарегистрированы!</h2>
  <p>Здравствуйте, %s %s!</p>
  <p>Вы зарегистрированы на мероприятие:</p>
  <table style="border-collapse: collapse; width: 100%%;">
    <tr><td style="padding: 8px; font-weight: bold;">Название</td><td style="padding: 8px;">%s</td></tr>
    <tr><td style="padding: 8px; font-weight: bold;">Дата</td><td style="padding: 8px;">%s — %s</td></tr>
    <tr><td style="padding: 8px; font-weight: bold;">Место</td><td style="padding: 8px;">%s</td></tr>
  </table>
  <p>Ваш QR-код для входа:</p>
  <img src="data:image/png;base64,%s" alt="QR-код" style="width: 200px; height: 200px;" />
  <p style="color: #888; font-size: 12px;">Покажите этот QR-код на входе. Не передавайте его третьим лицам.</p>
</body>
</html>`,
		participant.FirstName, participant.LastName,
		event.Title,
		utils.FormatDateTime(event.StartAt),
		utils.FormatDateTime(event.EndAt),
		event.Location,
		qrBase64,
	)

	return s.send(participant.Email, subject, body)
}

// send отправляет HTML письмо через SMTP
func (s *EmailService) send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	boundary := "boundary_qr_email"

	var msg bytes.Buffer
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.cfg.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: =?UTF-8?B?%s?=\r\n",
		base64.StdEncoding.EncodeToString([]byte(subject))))
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundary))
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString("\r\n")

	// Кодируем тело письма построчно по 76 символов (стандарт MIME)
	encoded := base64.StdEncoding.EncodeToString([]byte(htmlBody))
	for len(encoded) > 76 {
		msg.WriteString(encoded[:76] + "\r\n")
		encoded = encoded[76:]
	}
	msg.WriteString(encoded + "\r\n")

	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg.Bytes())
}

// SendBulkRegistrationEmails отправляет письма списку участников (для импорта из Excel)
func (s *EmailService) SendBulkRegistrationEmails(participants []models.Participant, event *models.Event) []error {
	var errs []error
	for i := range participants {
		if err := s.SendRegistrationEmail(&participants[i], event); err != nil {
			errs = append(errs, fmt.Errorf("не удалось отправить письмо на %s: %w",
				participants[i].Email, err))
		}
	}
	return errs
}

// sanitize нужен чтобы strings импорт не ругался
var _ = strings.TrimSpace
