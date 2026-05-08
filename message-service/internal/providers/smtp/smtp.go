package smtp

import (
	"fmt"
	"net/smtp"

	"message-service/config"
	"message-service/internal/models"
)

type Provider struct {
	cfg config.SMTPConfig
}

func NewSMTPProvider(cfg config.SMTPConfig) *Provider {
	return &Provider{
		cfg: cfg,
	}
}

func (p *Provider) Send(message models.Message) error {

	auth := smtp.PlainAuth(
		"",
		p.cfg.Email,
		p.cfg.Password,
		p.cfg.Host,
	)

	// Use CRLF and proper MIME headers so clients render HTML.
	body := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		p.cfg.Email,
		message.To,
		message.Subject,
		message.Content,
	)

	addr := p.cfg.Host + ":" + p.cfg.Port

	return smtp.SendMail(
		addr,
		auth,
		p.cfg.Email,
		[]string{message.To},
		[]byte(body),
	)
}