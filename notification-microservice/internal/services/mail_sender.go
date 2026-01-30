package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"time"
)

type MailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	user     string
	log      *slog.Logger
}

func NewMailSender(log *slog.Logger) *MailSender {
	return &MailSender{
		log:      log,
		from:     os.Getenv("SMTP_FROM"),
		password: os.Getenv("SMTP_PASS"),
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
	}
}

func (ms *MailSender) SendEmail(to, subject, body string) error {
	const op = "services.SendEmail"
	log := slog.With("op", op, "to", to, "subject", subject)

	if to == "" {
		log.Error("recipient email is empty")
		return fmt.Errorf("%s: recipient email is empty", op)
	}

	addr := ms.buildAddr()
	log.Info("builded addr", slog.String("addr", addr))
	auth := smtp.PlainAuth(
		"",
		ms.user,
		ms.password,
		ms.smtpHost,
	)

	log.Info("smtp config",
		slog.String("host", os.Getenv("SMTP_HOST")),
		slog.String("port", os.Getenv("SMTP_PORT")),
		slog.String("user", os.Getenv("SMTP_USER")),
		slog.String("from", os.Getenv("SMTP_FROM")),
	)

	msg := ms.buildMessage(to, subject, body)

	log.Info("msg is:", slog.String("msg", string(msg)))
	return smtp.SendMail(
		addr,
		auth,
		ms.from,
		[]string{to},
		msg,
	)
}

func (ms *MailSender) SendWithRetry(
	ctx context.Context,
	to, subject, body string,
) error {

	const (
		maxAttempts = 3
		delay       = 5 * time.Second
	)

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

		ms.log.Warn(
			"failed to send email, retrying",
			slog.Int("attempt", attempt),
			slog.Any("error", err),
		)

		lastErr = err
		time.Sleep(delay)
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
