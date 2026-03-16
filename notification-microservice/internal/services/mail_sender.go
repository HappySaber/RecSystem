// internal/services/mail.go
package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"time"
)

// internal/services/mail.go — добавить конструктор
func NewMailSenderWithSMTPAndDelay(log *slog.Logger, smtp SMTPSender, delay time.Duration) *MailSender {
	return &MailSender{
		log:        log,
		from:       "test@example.com",
		password:   "testpass",
		smtpHost:   "smtp.example.com",
		smtpPort:   "587",
		user:       "testuser",
		smtp:       smtp,
		retryDelay: delay, // добавляем поле
	}
}

// SMTPSender интерфейс для отправки — позволяет мокать в тестах
type SMTPSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// realSMTPSender — реальная реализация через net/smtp
type realSMTPSender struct{}

func (r *realSMTPSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

type MailSender struct {
	from       string
	password   string
	smtpHost   string
	smtpPort   string
	user       string
	log        *slog.Logger
	smtp       SMTPSender
	retryDelay time.Duration // для настройки задержки между попытками
}

// NewMailSender — для продакшна
func NewMailSender(log *slog.Logger) *MailSender {
	return &MailSender{
		log:      log,
		from:     os.Getenv("SMTP_FROM"),
		password: os.Getenv("SMTP_PASS"),
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		smtp:     &realSMTPSender{}, // реальный SMTP
	}
}

// NewMailSenderWithSMTP — для тестов, принимает мок
func NewMailSenderWithSMTP(log *slog.Logger, smtp SMTPSender) *MailSender {
	return &MailSender{
		log:      log,
		from:     "test@example.com",
		password: "testpass",
		smtpHost: "smtp.example.com",
		smtpPort: "587",
		user:     "testuser",
		smtp:     smtp,
	}
}

func (ms *MailSender) SendEmail(to, subject, body string) error {
	const op = "services.SendEmail"
	log := ms.log.With("op", op, "to", to, "subject", subject)

	if to == "" {
		log.Error("recipient email is empty")
		return fmt.Errorf("%s: recipient email is empty", op)
	}

	addr := ms.buildAddr()
	auth := smtp.PlainAuth("", ms.user, ms.password, ms.smtpHost)
	msg := ms.buildMessage(to, subject, body)

	// вызываем через интерфейс а не напрямую
	return ms.smtp.SendMail(addr, auth, ms.from, []string{to}, msg)
}

func (ms *MailSender) SendWithRetry(ctx context.Context, to, subject, body string) error {
	const maxAttempts = 3

	// если delay не задан — используем продакшн значение
	delay := ms.retryDelay
	if delay == 0 {
		delay = 5 * time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := ms.SendEmail(to, subject, body)
		if err == nil {
			return nil
		}

		ms.log.Warn("failed to send email, retrying",
			slog.Int("attempt", attempt),
			slog.Any("error", err),
		)
		lastErr = err

		if attempt < maxAttempts {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("send email failed after retries: %w", lastErr)
}

func (ms *MailSender) buildMessage(to, subject, body string) []byte {
	return []byte(
		"From: " + ms.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)
}

func (ms *MailSender) buildAddr() string {
	return fmt.Sprintf("%s:%s", ms.smtpHost, ms.smtpPort)
}
