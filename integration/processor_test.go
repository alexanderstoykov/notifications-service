//go:build integration

package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alexanderstoykov/notifications-service/internal/cron/processor"
	"github.com/alexanderstoykov/notifications-service/internal/service"
	"github.com/alexanderstoykov/notifications-service/internal/service/slack"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

func TestProcessorProcess(t *testing.T) {
	ctx := context.Background()

	processor := processor.NewProcessor(cfg.Cron, conn, notificationRepo)

	t.Run("when notifications batch processing succeeded", func(t *testing.T) {
		purgeTables(ctx)

		storedNotifications := []*storage.Notification{
			{
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 1",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, {
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 2",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, {
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 3",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		}

		err := insertNotifications(ctx, storedNotifications)
		assert.NoError(t, err)

		mockClient := service.NewMockSender(gomock.NewController(t))
		mockClient.EXPECT().Send(gomock.Any()).Times(cfg.Cron.BatchSize).Return(nil)
		processor.RegisterNotifier(storage.NotificationTypeSlack, slack.NewService(mockClient))

		err = processor.Process(ctx, storage.NotificationTypeSlack)
		assert.NoError(t, err)

		got, err := listAllNotifications(ctx)
		assert.NoError(t, err)

		exp := []*storage.Notification{
			{
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 3",
				Status:  storage.NotificationStatusPending,
			}, {
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 1",
				Status:  storage.SentNotificationStatusSent,
			}, {
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 2",
				Status:  storage.SentNotificationStatusSent,
			},
		}
		if diff := deep.Equal(exp, got); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("when notifications processing failed", func(t *testing.T) {
		purgeTables(ctx)

		storedNotifications := []*storage.Notification{
			{
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 1",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, {
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 2",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, {
				ID:        uuid.New(),
				Type:      storage.NotificationTypeSlack,
				Message:   "slack message 3",
				Status:    storage.NotificationStatusPending,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		}

		err := insertNotifications(ctx, storedNotifications)
		assert.NoError(t, err)

		mockClient := service.NewMockSender(gomock.NewController(t))
		mockClient.EXPECT().Send(gomock.Any()).Times(cfg.Cron.BatchSize).Return(errors.New("oops"))
		processor.RegisterNotifier(storage.NotificationTypeSlack, slack.NewService(mockClient))

		err = processor.Process(ctx, storage.NotificationTypeSlack)
		assert.NoError(t, err)

		got, err := listAllNotifications(ctx)
		assert.NoError(t, err)

		exp := []*storage.Notification{
			{
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 1",
				Status:  storage.NotificationStatusPending,
			},
			{
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 2",
				Status:  storage.NotificationStatusPending,
			}, {
				Type:    storage.NotificationTypeSlack,
				Message: "slack message 3",
				Status:  storage.NotificationStatusPending,
			},
		}
		if diff := deep.Equal(exp, got); diff != nil {
			t.Error(diff)
		}
	})
}
