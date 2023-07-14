package client

import (
	"github.com/alexanderstoykov/notifications-service/internal/service"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/alexanderstoykov/notifications-service/config"
)

type SMSClient struct {
	sender string
	client *twilio.RestClient
}

func NewSMSClient(config config.SMSConfig) *SMSClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.SID,
		Password: config.Token,
	})

	return &SMSClient{client: client, sender: config.Number}
}

func (c *SMSClient) Send(request *service.Message) error {
	params := &twilioApi.CreateMessageParams{}

	params.SetTo(request.Receiver)
	params.SetFrom(c.sender)
	params.SetBody(request.Message)

	if _, err := c.client.Api.CreateMessage(params); err != nil {
		return err
	}

	return nil
}
