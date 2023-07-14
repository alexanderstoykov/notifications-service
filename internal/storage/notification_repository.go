package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type NotificationRepository struct {
	conn *Connection
}

func NewNotificationRepository(conn *Connection) *NotificationRepository {
	return &NotificationRepository{conn: conn}
}

func (r *NotificationRepository) InsertOne(ctx context.Context, notification *Notification) error {
	notification.CreatedAt = time.Now().UTC()
	notification.UpdatedAt = time.Now().UTC()
	notification.ID = uuid.New()

	_, err := r.conn.DB(ctx).ExecContext(
		ctx, "INSERT INTO notifications"+
			"(id,type,receiver,message,status,created_at,updated_at) "+
			"VALUES "+
			"($1,$2,$3,$4,$5,$6,$7)",
		notification.ID,
		notification.Type,
		notification.Receiver,
		notification.Message,
		notification.Status,
		notification.CreatedAt,
		notification.UpdatedAt,
	)

	return err
}

func (r NotificationRepository) ListUnprocessedByTypeForUpdate(
	ctx context.Context,
	notificationType NotificationType,
	batchSize int,
) ([]*Notification, error) {
	notifications := make([]*Notification, 0, batchSize)

	err := r.conn.DB(ctx).
		SelectContext(ctx,
			&notifications,
			"SELECT * FROM notifications "+
				"WHERE type = $1 "+
				"AND status != 'SENT' "+
				"LIMIT $2 FOR UPDATE",
			notificationType,
			batchSize,
		)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r NotificationRepository) UpdateList(
	ctx context.Context,
	status NotificationStatus,
	notifications []*Notification,
) error {
	ids := make([]uuid.UUID, 0, len(notifications))
	for _, notification := range notifications {
		ids = append(ids, notification.ID)
	}

	query, args, err := sqlx.In("UPDATE notifications set status= ?, updated_at=? where id IN (?)",
		status,
		time.Now().UTC(),
		ids,
	)
	if err != nil {
		return err
	}

	query = r.conn.DB(ctx).Rebind(query)
	if _, err = r.conn.DB(ctx).ExecContext(ctx, query, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
