package services_test

import (
	"context"
	"errors"
	"log/slog"
	"net/smtp"
	"os"
	"testing"

	"notifications/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// мок SMTP sender
type mockSMTP struct {
	mock.Mock
}

func (m *mockSMTP) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}

func newTestMailSender(smtp services.SMTPSender) *services.MailSender {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	return services.NewMailSenderWithSMTP(log, smtp)
}

// SendEmail

func TestSendEmail_Success(t *testing.T) {
	smtpMock := new(mockSMTP)
	sender := newTestMailSender(smtpMock)

	// ожидаем вызов с любыми аргументами — возвращаем nil (успех)
	smtpMock.On("SendMail",
		mock.Anything,                // addr
		mock.Anything,                // auth
		mock.Anything,                // from
		[]string{"user@example.com"}, // to
		mock.Anything,                // msg
	).Return(nil).Once()

	err := sender.SendEmail("user@example.com", "Welcome!", "Hello!")

	require.NoError(t, err)
	smtpMock.AssertExpectations(t)
}

func TestSendEmail_EmptyRecipient(t *testing.T) {
	smtpMock := new(mockSMTP)
	sender := newTestMailSender(smtpMock)

	err := sender.SendEmail("", "Welcome!", "Hello!")

	require.Error(t, err)
	assert.ErrorContains(t, err, "recipient email is empty")
	// smtp не должен вызываться если получатель пустой
	smtpMock.AssertNotCalled(t, "SendMail")
}

func TestSendEmail_SMTPError(t *testing.T) {
	smtpMock := new(mockSMTP)
	sender := newTestMailSender(smtpMock)

	smtpMock.On("SendMail",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		[]string{"user@example.com"},
		mock.Anything,
	).Return(errors.New("connection refused")).Once()

	err := sender.SendEmail("user@example.com", "Welcome!", "Hello!")

	require.Error(t, err)
	assert.ErrorContains(t, err, "connection refused")
	smtpMock.AssertExpectations(t)
}

func TestSendEmail_MessageFormat(t *testing.T) {
	smtpMock := new(mockSMTP)
	sender := newTestMailSender(smtpMock)

	// проверяем что сообщение содержит нужные заголовки
	smtpMock.On("SendMail",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(msg []byte) bool {
			msgStr := string(msg)
			return assert.Contains(t, msgStr, "Subject: Test Subject") &&
				assert.Contains(t, msgStr, "To: user@example.com") &&
				assert.Contains(t, msgStr, "Content-Type: text/plain") &&
				assert.Contains(t, msgStr, "Test Body")
		}),
	).Return(nil).Once()

	err := sender.SendEmail("user@example.com", "Test Subject", "Test Body")

	require.NoError(t, err)
	smtpMock.AssertExpectations(t)
}

// SendWithRetry

func TestSendWithRetry_SuccessFirstAttempt(t *testing.T) {
	smtpMock := new(mockSMTP)
	sender := newTestMailSender(smtpMock)

	// успех с первой попытки — вызывается только раз
	smtpMock.On("SendMail",
		mock.Anything, mock.Anything, mock.Anything,
		[]string{"user@example.com"}, mock.Anything,
	).Return(nil).Once()

	err := sender.SendWithRetry(context.Background(), "user@example.com", "Welcome!", "Hello!")

	require.NoError(t, err)
	smtpMock.AssertExpectations(t)
}

func TestSendWithRetry_SuccessAfterRetry(t *testing.T) {
	smtpMock := new(mockSMTP)

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// используем кастомный sender с нулевым delay для быстрого теста
	sender := services.NewMailSenderWithSMTPAndDelay(log, smtpMock, 0)

	// первые два раза ошибка, третий — успех
	smtpMock.On("SendMail",
		mock.Anything, mock.Anything, mock.Anything,
		[]string{"user@example.com"}, mock.Anything,
	).Return(errors.New("temporary error")).Twice()

	smtpMock.On("SendMail",
		mock.Anything, mock.Anything, mock.Anything,
		[]string{"user@example.com"}, mock.Anything,
	).Return(nil).Once()

	err := sender.SendWithRetry(context.Background(), "user@example.com", "Welcome!", "Hello!")

	require.NoError(t, err)
	smtpMock.AssertExpectations(t)
}

func TestSendWithRetry_AllAttemptsFail(t *testing.T) {
	smtpMock := new(mockSMTP)

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sender := services.NewMailSenderWithSMTPAndDelay(log, smtpMock, 0)

	smtpMock.On("SendMail",
		mock.Anything, mock.Anything, mock.Anything,
		[]string{"user@example.com"}, mock.Anything,
	).Return(errors.New("smtp unavailable")).Times(3)

	err := sender.SendWithRetry(context.Background(), "user@example.com", "Welcome!", "Hello!")

	require.Error(t, err)
	assert.ErrorContains(t, err, "send email failed after retries")
	assert.ErrorContains(t, err, "smtp unavailable")
	smtpMock.AssertExpectations(t)
}

func TestSendWithRetry_ContextCancelled(t *testing.T) {
	smtpMock := new(mockSMTP)

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sender := services.NewMailSenderWithSMTPAndDelay(log, smtpMock, 0)

	// контекст уже отменён до начала
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sender.SendWithRetry(ctx, "user@example.com", "Welcome!", "Hello!")

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	// smtp не должен вызываться
	smtpMock.AssertNotCalled(t, "SendMail")
}
