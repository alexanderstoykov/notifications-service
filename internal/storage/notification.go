package storage

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeSMS   NotificationType = "SMS"
	NotificationTypeSlack NotificationType = "SLACK"
	NotificationTypeEmail NotificationType = "EMAIL"
)

type NotificationStatus string

const (
	NotificationStatusFailed   NotificationStatus = "FAILED"
	SentNotificationStatusSent NotificationStatus = "SENT"
	NotificationStatusPending  NotificationStatus = "PENDING"
)

type Notification struct {
	ID        uuid.UUID          `db:"id" deep:"-"`
	Type      NotificationType   `db:"type"`
	Receiver  string             `db:"receiver"`
	Message   string             `db:"message"`
	Status    NotificationStatus `db:"status"`
	CreatedAt time.Time          `db:"created_at" deep:"-"`
	UpdatedAt time.Time          `db:"updated_at" deep:"-"`
}
