// internal/repository/notification_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type notificationRepository struct {
	db *PostgresDB
}

func NewNotificationRepository(db *PostgresDB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(n *domain.Notification) error {
	query := `
        INSERT INTO notifications (user_id, title, content, notification_type, is_read)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING notification_id, created_at`

	return r.db.DB.QueryRow(
		query,
		n.UserID,
		n.Title,
		n.Content,
		n.NotificationType,
		n.IsRead,
	).Scan(&n.ID, &n.CreatedAt)
}

func (r *notificationRepository) GetByID(id int64) (*domain.Notification, error) {
	n := &domain.Notification{}
	query := `
        SELECT notification_id, user_id, title, content, notification_type, is_read, created_at
        FROM notifications 
        WHERE notification_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&n.ID,
		&n.UserID,
		&n.Title,
		&n.Content,
		&n.NotificationType,
		&n.IsRead,
		&n.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidOperation
	}
	return n, err
}

func (r *notificationRepository) GetByUserID(userID int64, limit, offset int) ([]*domain.Notification, error) {
	query := `
        SELECT notification_id, user_id, title, content, notification_type, is_read, created_at
        FROM notifications 
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Title,
			&n.Content,
			&n.NotificationType,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *notificationRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	query := `
        SELECT COUNT(*)
        FROM notifications
        WHERE user_id = $1 AND is_read = false`

	err := r.db.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *notificationRepository) MarkAsRead(id int64, userID int64) error {
	query := `
        UPDATE notifications
        SET is_read = true
        WHERE notification_id = $1 AND user_id = $2`

	result, err := r.db.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrInvalidOperation
	}

	return nil
}

func (r *notificationRepository) MarkAllAsRead(userID int64) error {
	query := `
        UPDATE notifications
        SET is_read = true
        WHERE user_id = $1 AND is_read = false`

	_, err := r.db.DB.Exec(query, userID)
	return err
}
