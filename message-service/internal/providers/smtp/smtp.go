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

	body := fmt.Sprintf(
		"Subject: %s\n\n%s",
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