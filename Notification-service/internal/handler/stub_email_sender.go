package handler

import "go.uber.org/zap"

type EmailSender interface {
	Send(to, subject, body string) error
}

type StubEmailSender struct {
	log *zap.Logger
}

func NewStubEmailSender(log *zap.Logger) *StubEmailSender {
	return &StubEmailSender{log: log}
}

func (s *StubEmailSender) Send(to, subject, body string) error {
	s.log.Info("FAKE EMAIL SENT", zap.String("to", to), zap.String("subject", subject))
	return nil
}
