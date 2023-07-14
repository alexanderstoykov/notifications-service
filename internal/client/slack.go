package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexanderstoykov/notifications-service/config"
	"github.com/alexanderstoykov/notifications-service/internal/service"
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type SlackClient struct {
	url    string
	client HttpDoer
}

func NewSlackClient(config config.SlackConfig) *SlackClient {
	return &SlackClient{
		url:    config.Webhook,
		client: &http.Client{},
	}
}

func (c *SlackClient) Send(message *service.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("could not close body: %s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}
