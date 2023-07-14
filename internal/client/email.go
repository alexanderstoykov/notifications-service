package client

import (
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/alexanderstoykov/notifications-service/config"
	"github.com/alexanderstoykov/notifications-service/internal/service"
)

type EmailClient struct {
	client *gomail.Dialer
}

func NewEmailClient(config config.MailConfig) *EmailClient {
	client := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	return &EmailClient{client: client}
}

func (c *EmailClient) Send(message *service.Message) error {
	m := gomail.NewMessage()
	m.SetHeader("From", message.Sender)
	m.SetHeader("To", message.Receiver)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text/plain", message.Message)

	err := c.client.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
