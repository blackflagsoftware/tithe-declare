package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/blackflagsoftware/tithe-declare/config"
	"github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
)

//go:generate mockgen -source=email.go -destination=mock.go -package=email
type (
	Emailer interface {
		SendReset(context.Context, string, string) error
		SendReminder(context.Context, []string, string) error
	}

	Email struct{}
)

func EmailInit() Emailer {
	if config.E.Host != "" {
		return &Email{}
	}
	return &MockEmailer{}
}

func (e Email) SendReset(ctx context.Context, toEmail, resetToken string) error {
	from := config.E.From
	pwd := config.E.Pwd
	host := config.E.Host
	port := config.E.GetEmailPort()

	auth := smtp.PlainAuth("", from, pwd, host)
	to := []string{toEmail}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Reset Password Instructions\r\n\r\nTo reset your password: %s?email=%s&token=%s\r\n", toEmail, config.E.ResetUrl, toEmail, resetToken))
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, to, msg); err != nil {
		logging.Default.Println("unable to send email:", err)
		return err
	}
	return nil
}

func (e Email) SendReminder(ctx context.Context, toEmail []string, body string) error {
	from := config.E.From
	pwd := config.E.Pwd
	host := config.E.Host
	port := config.E.GetEmailPort()

	auth := smtp.PlainAuth("", from, pwd, host)
	to := toEmail
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Upcoming Tithing Declarations\r\n\r\n%s\r\n", toEmail, body))
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, to, msg); err != nil {
		logging.Default.Println("unable to send email for SendReminder:", err)
		return err
	}
	return nil
}
