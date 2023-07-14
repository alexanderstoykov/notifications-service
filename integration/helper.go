//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/alexanderstoykov/notifications-service/config"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

var (
	cfg              *config.Config
	conn             *storage.Connection
	notificationRepo *storage.NotificationRepository
)

func init() {
	var err error

	if err = godotenv.Load("../.env.test"); err != nil {
		log.Fatal(err)
	}

	cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err = storage.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("error while connecting db: %s", err)
	}

	notificationRepo = storage.NewNotificationRepository(conn)
}

func listAllNotifications(ctx context.Context) ([]*storage.Notification, error) {
	var notifications []*storage.Notification

	err := conn.DB(ctx).
		SelectContext(ctx,
			&notifications,
			"SELECT * FROM notifications",
		)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func insertNotifications(ctx context.Context, notifications []*storage.Notification) error {
	for _, notification := range notifications {
		if err := notificationRepo.InsertOne(ctx, notification); err != nil {
			return err
		}
	}

	return nil
}

func purgeTables(ctx context.Context) {
	_, err := conn.DB(ctx).ExecContext(ctx, "TRUNCATE TABLE notifications;")
	if err != nil {
		panic(fmt.Sprintf("an error occurred cleaning the feed table for tests: %s", err))
	}
}
