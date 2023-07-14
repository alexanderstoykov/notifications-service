package processor

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/alexanderstoykov/notifications-service/config"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

type NotificationsRepository interface {
	ListUnprocessedByTypeForUpdate(context.Context, storage.NotificationType, int) ([]*storage.Notification, error)
	UpdateList(context.Context, storage.NotificationStatus, []*storage.Notification) error
}

type Notifier interface {
	Notify(ctx context.Context, notification *storage.Notification) error
}

type Transactioner interface {
	Tx(ctx context.Context, txFunc storage.TxFunc) error
}

type Processor struct {
	config                  config.CronConfig
	tx                      Transactioner
	notificationsRepository NotificationsRepository
	notifiers               map[storage.NotificationType]Notifier
}

func NewProcessor(
	config config.CronConfig,
	tx Transactioner,
	repository NotificationsRepository,
) *Processor {
	return &Processor{
		config:                  config,
		tx:                      tx,
		notificationsRepository: repository,
	}
}

func (p *Processor) Process(ctx context.Context, notificationType storage.NotificationType) error {
	return p.tx.Tx(ctx, func(ctx context.Context) error {
		notifications, err := p.notificationsRepository.ListUnprocessedByTypeForUpdate(
			ctx,
			notificationType,
			p.config.BatchSize,
		)
		if err != nil {
			return errors.WithStack(err)
		}

		if len(notifications) == 0 {
			return nil
		}

		log.Printf("processing %d notification", len(notifications))

		sentNotifications := make([]*storage.Notification, 0, len(notifications))
		for _, notification := range notifications {
			// nolint: govet
			notifier, err := p.resolveNotifier(notificationType)
			if err != nil {
				return errors.WithStack(err)
			}

			if err = notifier.Notify(ctx, notification); err != nil {
				log.Printf("could not process notification for type %s: %s", notificationType, err)

				continue
			}

			notification.Status = storage.SentNotificationStatusSent

			sentNotifications = append(sentNotifications, notification)
		}

		if len(sentNotifications) > 0 {
			log.Printf("%d notifications processed", len(sentNotifications))

			err = p.notificationsRepository.UpdateList(ctx, storage.SentNotificationStatusSent, notifications)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

func (p *Processor) RegisterNotifier(notificationType storage.NotificationType, notifier Notifier) {
	if p.notifiers == nil {
		p.notifiers = make(map[storage.NotificationType]Notifier)
	}

	log.Printf("registering %s notifier", notificationType)

	p.notifiers[notificationType] = notifier
}

func (p *Processor) resolveNotifier(notificationType storage.NotificationType) (Notifier, error) {
	notifier, ok := p.notifiers[notificationType]
	if !ok {
		return nil, newNotifierNotFoundError(notificationType)
	}

	return notifier, nil
}
