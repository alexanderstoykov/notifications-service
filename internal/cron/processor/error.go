package processor

import (
	"fmt"

	"github.com/alexanderstoykov/notifications-service/internal/storage"
)

type NotifierNotFoundError struct {
	notifierType storage.NotificationType
}

func newNotifierNotFoundError(notifierType storage.NotificationType) *NotifierNotFoundError {
	return &NotifierNotFoundError{notifierType: notifierType}
}

func (e *NotifierNotFoundError) Error() string {
	return fmt.Sprintf("notifier for %s notifications not found", e.notifierType)
}
