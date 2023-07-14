package main

import (
	"context"
	"flag"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/alexanderstoykov/notifications-service/config"
	"github.com/alexanderstoykov/notifications-service/internal/client"
	"github.com/alexanderstoykov/notifications-service/internal/cron/processor"
	"github.com/alexanderstoykov/notifications-service/internal/service/email"
	"github.com/alexanderstoykov/notifications-service/internal/service/slack"
	"github.com/alexanderstoykov/notifications-service/internal/service/sms"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		cancel()
	}()

	conn, err := storage.NewConnection(cfg.Database)
	if err != nil {
		return err
	}

	notificationRepo := storage.NewNotificationRepository(conn)
	processor := processor.NewProcessor(cfg.Cron, conn, notificationRepo)

	notificationTypeArg := flag.String("type", "slack", "notification type to process")

	flag.Parse()

	var notificationType storage.NotificationType

	switch *notificationTypeArg {
	case "email":
		notificationType = storage.NotificationTypeEmail

		emailClient := client.NewEmailClient(cfg.Email)
		processor.RegisterNotifier(notificationType, email.NewService(emailClient))
	case "sms":
		notificationType = storage.NotificationTypeSMS

		smsClient := client.NewSMSClient(cfg.SMS)
		processor.RegisterNotifier(notificationType, sms.NewService(smsClient))
	case "slack":
		notificationType = storage.NotificationTypeSlack

		slackClient := client.NewSlackClient(cfg.Slack)
		processor.RegisterNotifier(notificationType, slack.NewService(slackClient))
	}

	log.Printf("%s cron started", notificationType)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := processor.Process(ctx, notificationType); err != nil {
				return errors.WithStack(err)
			}

			time.Sleep(cfg.Cron.Interval)
		}
	}
}
