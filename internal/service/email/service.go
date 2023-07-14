package email

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alexanderstoykov/notifications-service/internal/service"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

type Service struct {
	client service.Sender
	sender string
}

func NewService(client service.Sender) *Service {
	return &Service{client: client}
}

func (s Service) Notify(ctx context.Context, notification *storage.Notification) error {
	message := &service.Message{
		Sender:   s.sender,
		Message:  notification.Message,
		Receiver: notification.Receiver,
	}

	err := s.client.Send(message)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
