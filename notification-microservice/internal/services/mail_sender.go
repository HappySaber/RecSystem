package services

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
)

type MailSender struct {
	from     string
	to       string
	password string
	smtpHost string
	smtpPort string
	user     string
	log      *slog.Logger
}

func New(log *slog.Logger) *MailSender {
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
	auth := smtp.PlainAuth(
		"",
		ms.user,
		ms.password,
		ms.smtpHost,
	)
	msg := []byte(
		"From: " + ms.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)
	return smtp.SendMail(
		addr,
		auth,
		ms.from,
		[]string{to},
		msg,
	)
}

func (ms *MailSender) buildAddr() string {
	return fmt.Sprintf("%s:%s", ms.smtpHost, ms.smtpPort)
}
