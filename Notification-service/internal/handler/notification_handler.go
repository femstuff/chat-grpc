package handler

import (
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

type EmailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	log      *zap.Logger
}

func NewEmailSender(from, pass, smtpHost, smtpPort string, log *zap.Logger) *EmailSender {
	return &EmailSender{
		from:     from,
		password: pass,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		log:      log,
	}
}

func (s *EmailSender) Send(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", s.from, to, subject, body)
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		s.log.Error("Failed to send email", zap.Error(err))
		return err
	}

	s.log.Info("email sent", zap.String("to", to))
	return nil
}
