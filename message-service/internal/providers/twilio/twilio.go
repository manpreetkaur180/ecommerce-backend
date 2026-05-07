package twilio

import (
	"message-service/config"
	"message-service/internal/models"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/twilio/twilio-go"
)

type Provider struct {
	cfg config.TwilioConfig
}

func NewTwilioProvider(
	cfg config.TwilioConfig,
) *Provider {

	return &Provider{
		cfg: cfg,
	}
}

func (p *Provider) Send(
	message models.Message,
) error {

	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: p.cfg.AccountSID,
			Password: p.cfg.AuthToken,
		},
	)

	params := &openapi.CreateMessageParams{}

	params.SetTo(message.To)
	params.SetFrom(p.cfg.FromNumber)
	params.SetBody(message.Content)

	_, err := client.Api.CreateMessage(params)

	return err
}